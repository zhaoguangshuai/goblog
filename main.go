package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func homeHandler(w http.ResponseWriter, r *http.Request)  {
	//w.Header().Set("Content-Type","text/html;charset=utf-8")
	fmt.Fprint(w,"<h1>hello，欢迎来到goblog！</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "text/html;charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter,r *http.Request)  {
	//设置返回行的内容类型是xml，还是json，还是文本内容；该设置为html类型，不然解析不了
	//w.Header().Set("Content-Type","text/html;charset=utf-8")
	//设置返回的http状态码为404
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w,"<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func articlesShowHandler(w http.ResponseWriter,r *http.Request)  {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "文章 ID:"+id)
}

func articlesIndexHandler(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprint(w,"访问文章列表11")
}

func articlesStoreHandler(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprint(w,"创建新的文章")
}

//添加一个设置返回内容格式的中间件
func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//1.设置标头
		w.Header().Set("Content-Type","text/html;charset=utf-8")
		//2.继续处理请求
		next.ServeHTTP(w,r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//1.除首页以外，移除所有请求路径后面的斜杠
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path,"/")
		}
		//2.将请求继续传递下去
		next.ServeHTTP(w,r)
	})
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")

	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}",articlesShowHandler).Methods("GET").Name("articles.show")

	router.HandleFunc("/articles",articlesIndexHandler).Methods("GET").Name("articles.index")

	router.HandleFunc("/articles",articlesStoreHandler).Methods("POST").Name("articles.shore")

	//自定义404页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	//中间件：强制内容类型为HTML
	router.Use(forceHTMLMiddleware)

	//通过命名路由获取 URL 事例
	homeURL,_ := router.Get("home").URL()
	fmt.Println(homeURL)
	articleURL,_ := router.Get("articles.show").URL("id","23")
	fmt.Println(articleURL)


	http.ListenAndServe(":3000", removeTrailingSlash(router))

}
