package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"goblog/bootstrap"
	"goblog/pkg/database"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"net/http"
	"strconv"
	"strings"
)

var router *mux.Router
var db *sql.DB

//Article 对应一条文章数据
type Article struct {
	Title,Body  string
	ID			int64
}

func getArticleByID(id string) (Article,error) {
	article := Article{}
	query := "select * from articles where id = ?"
	err := db.QueryRow(query,id).Scan(&article.ID,&article.Title,&article.Body)
	return article,err
}

//ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	Id		int64
	Errors map[string]string
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


	router.HandleFunc("/articles/{id:[0-9]+}/delete",articlesDeleteHandler).Methods("POST").Name("articles.delete")

	//中间件：强制内容类型为HTML
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))

}
