package bootstrap

import (
	"github.com/gorilla/mux"
	"goblog/pkg/route"
	"goblog/routes"
)

//为了解决循环引用才建立的
func SetupRoute() *mux.Router {
	//获取路由包的结构体指针对象
	router := mux.NewRouter()
	//注册路由
	routes.RegisterWebRoutes(router)
	//初始化路由
	route.SetRoute(router)//自定义的一些路由函数，里面需要用到路由包的结构体指针对象，所以在这里获取到传过去
	return router
}
