package view

import (
	"goblog/pkg/auth"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

//D 是  map[string]interface{}
type D map[string]interface{}

//Render 渲染通用视图
func Render(w io.Writer, data D,tplFiles ...string)  {
	RenderTemplate(w,"app",data,tplFiles...)
}

//RenderSimple 渲染简单的视图
func RenderSimple(w io.Writer, data D,tplFiles ...string)  {
	RenderTemplate(w,"simple",data,tplFiles...)
}

//RenderTemplate 渲染视图
func RenderTemplate(w io.Writer,name string, data D,tplFiles ...string)  {
	//1. 通用模版数据
	data["isLogined"] = auth.Check()

	//2.生成模版文件
	allFiles := getTemplateFiles(tplFiles...)

	//3. 解析所有模版文件
	tmpl,err := template.New("").Funcs(template.FuncMap{
		"RouteName2URL" : route.Name2URL,
	}).ParseFiles(allFiles...)
	logger.LogError(err)

	//4 渲染模版
	tmpl.ExecuteTemplate(w,name,data)
}

func getTemplateFiles(tplFiles ...string) []string {
	//1 设置模版的相对路径
	viewDir := "resources/views/"

	//2. 遍历传参文件列表 slice，设置正确的路径，支持 dir.filename 语法糖
	for i,f := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(f,".","/",-1)+".gohtml"
	}

	//3 所有布局模版文件 slice
	layoutFiles,err := filepath.Glob(viewDir+"layouts/*.gohtml")
	logger.LogError(err)

	//4 在slice 里新增我们的目标文件
	return append(layoutFiles,tplFiles...)
}
