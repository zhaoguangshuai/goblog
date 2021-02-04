package route

import (
	"github.com/gorilla/mux"
	"net/http"
)

//router 定义路由对象，并没有赋值
var route *mux.Router

//初始化路由对象 ， 可以理解为赋值
func SetRoute(r *mux.Router)  {
	route = r
}

// RouteName2URL 通过路由名称来获取 URL
func Name2URL(routeName string, pairs ...string) string {
	url, err := route.Get(routeName).URL(pairs...)
	if err != nil {
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