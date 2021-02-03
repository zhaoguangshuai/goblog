package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var router = mux.NewRouter()
var db *sql.DB

func homeHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type","text/html;charset=utf-8")
	fmt.Fprint(w, "<h1>hello，欢迎来到goblog！</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "text/html;charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	//设置返回行的内容类型是xml，还是json，还是文本内容；该设置为html类型，不然解析不了
	//w.Header().Set("Content-Type","text/html;charset=utf-8")
	//设置返回的http状态码为404
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	//1.获取URL参数
	id := getRouteVariable("id",r)

	//2.读取对应的文章数据
	article := Article{}
	article,err := getArticleByID(id)

	//3.如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			//3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w,"404 文章未找到")
		} else {
			//3.2 数据库错误
			checkError(err)//记录错误日志
			//设置返回的http状态码
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"500 服务器内部错误")
		}
	} else {
		//4.读取成功
		tmpl,err := template.ParseFiles("resources/views/articles/show.gohtml")
		checkError(err)
		tmpl.Execute(w,article)
	}

}

//Article 对应一条文章数据
type Article struct {
	Title,Body  string
	ID			int64
}

func articlesEditHandler(w http.ResponseWriter,r *http.Request)  {
	//1.获取url参数
	id := getRouteVariable("id",r)

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
			checkError(err)
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
		checkError(err)

		tmpl.Execute(w,data)
	}
}

func articlesUpdateHandler(w http.ResponseWriter,r *http.Request)  {
	//1.获取URL参数
	id := getRouteVariable("id",r)

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
			checkError(err)
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
				checkError(err)
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
			checkError(err)
			tmpl.Execute(w,data)
		}
	}

}
/*
通过传参 URL 路由参数名称获取值
*/
func getRouteVariable(parameterName string,r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

func getArticleByID(id string) (Article,error) {
	article := Article{}
	query := "select * from articles where id = ?"
	err := db.QueryRow(query,id).Scan(&article.ID,&article.Title,&article.Body)
	return article,err
}


func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "访问文章列表11")
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
			checkError(err)
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

func initDB()  {
	var err error
	config := mysql.Config{
		User:			"root",
		Passwd:			"zgs8653406",
		Addr:			"127.0.0.1:3306",
		Net: 			"tcp",
		DBName: 		"goblog",
		AllowNativePasswords: true,
	}

	//准备连接数据库
	db,err = sql.Open("mysql",config.FormatDSN())
	checkError(err)
	//fmt.Println(config.FormatDSN())
	//fmt.Println(db)

	//设置最大连接数
	db.SetMaxOpenConns(25)
	//设置最大空闲连接数
	db.SetMaxIdleConns(25)
	//设置每个连接的过期时间
	db.SetConnMaxLifetime(5 * time.Minute)

	//尝试连接，失败会报错
	//err = db.Ping()
	//checkError(err)

}

func checkError(err error)  {
	if err != nil {
		log.Fatal(err)
	}
}

func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
); `

	_, err := db.Exec(createArticlesSQL)
	checkError(err)
}

func main() {
	initDB()
	createTables()

	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")

	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")

	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")

	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.shore")

	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	router.HandleFunc("/articles/{id:[0-9]+}/edit",articlesEditHandler).Methods("GET").Name("articles.edit")

	router.HandleFunc("/articles/{id:[0-9]+}",articlesUpdateHandler).Methods("POST").Name("articles.update")

	//自定义404页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	//中间件：强制内容类型为HTML
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))

}
