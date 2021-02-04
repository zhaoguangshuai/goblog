package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strconv"
	"unicode/utf8"
)

// ArticlesController 文章相关页面
type ArticlesController struct {
}
//Show 文章详情页
func (* ArticlesController) Show(w http.ResponseWriter, r *http.Request)  {
	//1.获取URL参数
	id := route.GetRouteVariable("id",r)

	//2.读取对应的文章数据
	article,err := article.Get(id)

	//3.如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w,"404 文章未找到")
		} else {
			//3.2 数据库错误
			logger.LogError(err)//记录错误日志
			//设置返回的http状态码
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"500 服务器内部错误")
		}
	} else {
		//4.读取成功
		//tmpl,err := template.ParseFiles("resources/views/articles/show.gohtml")
		tmpl,err := template.New("show.gohtml").Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
			"Int64ToString": types.Int64ToString,
		}).ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)
		tmpl.Execute(w,article)
	}
}

//首页文章列表
func (* ArticlesController) Index(w http.ResponseWriter, r *http.Request) {

	//1. 获取结果集
	articles,err := article.GetAll()

	if err != nil {
		//数据库错误
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w,"500 服务器内部错误")
	} else {
		//2.加载模版
		tmpl,err := template.ParseFiles("resources/views/articles/index.gohtml")
		logger.LogError(err)

		//3.渲染模版，将所有文章的数据传输进去
		tmpl.Execute(w,articles)
	}

}

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title,Body		string
	URL				string
	Errors			map[string]string
}

func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request)  {
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

func validateArticleFormData(title string, body string) map[string]string {
	errors := make(map[string]string)
	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	return errors
}

// Store 文章创建页面
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)

	// 检查是否有错误
	if len(errors) == 0 {
		_article := article.Article{
			Title: title,
			Body:  body,
		}
		_article.Create()
		if _article.ID > 0 {
			fmt.Fprint(w, "插入成功，ID 为"+strconv.FormatInt(_article.ID, 10))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章失败，请联系管理员")
		}
	} else {
		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")

		logger.LogError(err)

		tmpl.Execute(w, data)
	}
}

