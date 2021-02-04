package route

import "github.com/gorilla/mux"

//router 路由对象
var Router *mux.Router

// initialize 初始化路由
func Initialize()  {
	Router = mux.NewRouter()
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