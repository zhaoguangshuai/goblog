package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goblog/bootstrap"
	"goblog/pkg/database"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

var router = route.Router
var db *sql.DB

//Article 对应一条文章数据
type Article struct {
	Title,Body  string
	ID			int64
}

//Link 方法用来生成文章链接
func (a Article) Link() string {
	showURL,err := router.Get("articles.show").URL("id",strconv.FormatInt(a.ID,10))
	if err != nil {
		logger.LogError(err)
		return  ""
	}
	return showURL.String()
}

func articlesEditHandler(w http.ResponseWriter,r *http.Request)  {
	//1.获取url参数
	id := route.GetRouteVariable("id",r)

	//2.读取对应的文章数据
	article,err := getArticleByID(id)

	//3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			//3.1 数据为找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w,"404 文章未找到")
		} else {
			//3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"500 服务器内部错误")
		}
	} else {
		//4 读取成功，显示表单
		//updateURL,_ := router.Get("articles.update").URL("id",id)
		//将字符串转化为int64
		id1,_ := strconv.ParseInt(id,10,64)
		data := ArticlesFormData{
			Title: 			article.Title,
			Body: 			article.Body,
			Id: 			id1,
			Errors: 		nil,
		}
		tmpl,err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logger.LogError(err)

		tmpl.Execute(w,data)
	}
}

func articlesUpdateHandler(w http.ResponseWriter,r *http.Request)  {
	//1.获取URL参数
	id := route.GetRouteVariable("id",r)

	//获取对应的文章数据
	_,err := getArticleByID(id)

	//3.如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			//3.1数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w,"404 文章未找到")
		} else {
			//3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"500 服务器内部错误")
		}
	} else {
		//4. 未出现错误

		//4.1 表单验证
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := make(map[string]string)

		// 验证title
		if title == "" {
			errors["title"] = "标题不能为空"
		} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
			errors["title"] = "标题长度介于 3-40"
		}

		//验证内容
		if body == "" {
			errors["body"] = "内容不能为空"
		} else if utf8.RuneCountInString(body) < 10 {
			errors["body"] = "内容长度需大于或等于 10 个字节"
		}

		if len(errors) == 0 {
			//4.2 表单验证通过，更新数据
			query := "update articles set title = ?,body = ? where id = ?"
			rs,err := db.Exec(query,title,body,id)
			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w,"500 服务器内部错误")
			}

			//更新成功，跳转到文章详情页
			if n,_ := rs.RowsAffected();n > 0 {
				http.Redirect(w,r,"/articles/"+id,http.StatusFound)
			} else {
				fmt.Fprint(w,"您没有做任何更改!")
			}
		} else {
			//4.3 表单验证不通过，显示理由
			id1,_ := strconv.ParseInt(id,10,64)
			data := ArticlesFormData{
				Title: title,
				Body: body,
				Id: id1,
				Errors: errors,
			}
			tmpl,err := template.ParseFiles("resources/views/articles/edit.gohtml")
			logger.LogError(err)
			tmpl.Execute(w,data)
		}
	}

}


func getArticleByID(id string) (Article,error) {
	article := Article{}
	query := "select * from articles where id = ?"
	err := db.QueryRow(query,id).Scan(&article.ID,&article.Title,&article.Body)
	return article,err
}


func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	//1. 执行查询语句，返回一个结果集
	rows,err := db.Query("select * from articles")
	logger.LogError(err)
	defer rows.Close()

	var articles []Article
	//2.循环读取结果
	for rows.Next() {
		var article Article
		//2.1 扫码每一行的结果并赋值到一个 article 对象中
		err := rows.Scan(&article.ID, &article.Title,&article.Body)
		logger.LogError(err)
		//2.2 将article追加到articles 的这个数组中
		articles = append(articles,article)
	}
	//2.3 检查遍历时是否发生错误
	err = rows.Err()
	logger.LogError(err)

	//3.加载模版
	tmpl,err := template.ParseFiles("resources/views/articles/index.gohtml")
	logger.LogError(err)
	//4.渲染模版，将所有文章的数据传输进去
	tmpl.Execute(w,articles)

}

//ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	Id		int64
	Errors map[string]string
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	errors := make(map[string]string)

	//验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	//验证内容
	if body == "" {
		errors["body"] = "标题不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	//检查是否有错误
	if len(errors) == 0 {
		lastInsertID,err := saveArticleToDB(title,body)
		if lastInsertID > 0 {
			fmt.Fprint(w, "插入成功，id为"+strconv.FormatInt(lastInsertID,10))
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"500 服务器内部错误")
		}

	} else {
		data := ArticlesFormData{
			Title: title,
			Body: body,
			Errors: errors,
		}

		tmpl,err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}
		tmpl.Execute(w,data)
	}

}

func saveArticleToDB(title string,body string) (int64,error) {
	//变量初始化
	var (
		id 		int64
		err		error
		rs		sql.Result
		stmt	*sql.Stmt
	)

	//1.获取一个 prepare 声明语句
	stmt,err = db.Prepare("insert into articles (title,body) values (?,?)")
	//例行的错误检测
	if err != nil {
		return 0, err
	}

	//2.在此函数运行结束后关闭此语句，防止占用sql连接
	defer stmt.Close()

	//3.执行请求，传参进入绑定的内容
	rs,err = stmt.Exec(title,body)//正式执行sql语句
	if err != nil {
		return 0, err
	}

	//4. 插入成功的花，会返回自增 id
	if id,err = rs.LastInsertId(); id > 0 {
		return id,nil
	}
	return 0,nil

}

//添加一个设置返回内容格式的中间件
func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//1.设置标头
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		//2.继续处理请求
		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//1.除首页以外，移除所有请求路径后面的斜杠
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		//2.将请求继续传递下去
		next.ServeHTTP(w, r)
	})
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		Errors: nil,
	}
	tmpl,err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w,data)
}

func articlesDeleteHandler(w http.ResponseWriter,r *http.Request)  {
	//1. 获取URL 参数
	id := route.GetRouteVariable("id",r)

	//2.读取对应的文章数据
	article,err := getArticleByID(id)

	//3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			//3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			//3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"500 服务器内部错误")
		}
	} else {
		//4. 未出现错误，执行删除操作
		rowsAffected, err := article.Delete()

		//4.1 发生错误
		if err != nil {
			//因该是 sql 报错了
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"500 服务器内部错误")
		} else {
			//4.2 未发生错误
			if rowsAffected > 0 {
				//重定向到文章列表页
				indexURL,_ := router.Get("articles.index").URL()
				http.Redirect(w,r,indexURL.String(),http.StatusFound)
			} else {
				//Edga case
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w,"404 文章未找到")
			}
		}
	}
}

//Delete 方法用以从数据库中删除单条记录
func (a Article) Delete() (rowsAffected int64,err error) {
	rs,err := db.Exec("delete from articles where id = "+strconv.FormatInt(a.ID,10))

	if err != nil {
		return 0,err
	}

	//更新成功，跳转到文章详情页
	if n,_ := rs.RowsAffected(); n > 0 {
		return n,nil
	}
	return 0,nil
}

func main() {
	database.Initialize()
	db = database.DB

	//route.Initialize()
	//router = route.Router
	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()

	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")

	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.shore")

	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	router.HandleFunc("/articles/{id:[0-9]+}/edit",articlesEditHandler).Methods("GET").Name("articles.edit")

	router.HandleFunc("/articles/{id:[0-9]+}",articlesUpdateHandler).Methods("POST").Name("articles.update")

	router.HandleFunc("/articles/{id:[0-9]+}/delete",articlesDeleteHandler).Methods("POST").Name("articles.delete")

	//中间件：强制内容类型为HTML
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))

}
