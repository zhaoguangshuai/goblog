package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type","text/html;charset=utf-8")
	fmt.Fprint(w,"<h1>hello，欢迎来到goblog！</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter,r *http.Request)  {
	//设置返回行的内容类型是xml，还是json，还是文本内容；该设置为html类型，不然解析不了
	w.Header().Set("Content-Type","text/html;charset=utf-8")
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
	fmt.Fprint(w,"访问文章列表")
}

func articlesStoreHandler(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprint(w,"创建新的文章")
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

	//通过命名路由获取 URL 事例
	homeURL,_ := router.Get("home").URL()
	fmt.Println(homeURL)
	articleURL,_ := router.Get("articles.show").URL("id","23")
	fmt.Println(articleURL)


	http.ListenAndServe(":3000", router)
}
