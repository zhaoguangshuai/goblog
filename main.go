package main

import (
	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	"net/http"

)

func main() {
	//初始化数据库和 ORM
	bootstrap.SetupDB()
	//注册路由
	router := bootstrap.SetupRoute()
	//服务启动，端口监听，并且除首页以外，移除所有请求路径后面的斜杠
	http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))

}
