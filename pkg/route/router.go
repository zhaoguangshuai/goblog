package route

import (
	"github.com/gorilla/mux"
	"goblog/routes"
	"net/http"
)

//router 路由对象
var Router *mux.Router

// initialize 初始化路由
func Initialize()  {
	Router = mux.NewRouter()
	routes.RegisterWebRoutes(Router)
}

// RouteName2URL 通过路由名称来获取 URL
func Name2URL(routeName string, pairs ...string) string {
	url, err := Router.Get(routeName).URL(pairs...)
	if err != nil {
		//checkError(err)
		return ""
	}

	return url.String()
}

/*
通过传参 URL 路由参数名称获取值
*/
func GetRouteVariable(parameterName string,r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}