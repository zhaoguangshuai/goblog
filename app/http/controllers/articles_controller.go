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
	"path/filepath"
	"strconv"
	"unicode/utf8"
)

// ArticlesController 文章相关页面
type ArticlesController struct {
}

//Show 文章详情页
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	//1.获取URL参数
	id := route.GetRouteVariable("id", r)

	//2.读取对应的文章数据
	article, err := article.Get(id)

	//3.如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			//3.2 数据库错误
			logger.LogError(err) //记录错误日志
			//设置返回的http状态码
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		//4.读取成功
		//tmpl,err := template.ParseFiles("resources/views/articles/show.gohtml")
		tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
			"Int64ToString": types.Int64ToString,
		}).ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)
		tmpl.Execute(w, article)
	}
}

//首页文章列表
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {

	//1. 获取结果集
	articles, err := article.GetAll()

	if err != nil {
		//数据库错误
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		//2.加载模版
		//tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
		//logger.LogError(err)
		//
		////3.渲染模版，将所有文章的数据传输进去
		//tmpl.Execute(w, articles)

		//2.0 设置模版相对路径
		viewDir := "resources/views"

		//2.1 所有布局模版文件 slice
		files,err := filepath.Glob(viewDir + "/layouts/*.gohtml")
		logger.LogError(err)

		//2.2 在slice 里新增我们的目标文件
		newFiles := append(files,viewDir+"/articles/index.gohtml")

		//2.3 解析模版文件
		tmpl,err := template.ParseFiles(newFiles...)
		logger.LogError(err)

		//2.4 渲染模版，将所有文章的数据传输进去
		tmpl.ExecuteTemplate(w,"app",articles)

	}

}

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL         string
	Id          int64
	Errors      map[string]string
}

func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		Errors: nil,
	}
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, data)
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

// Edit 文章更新页面
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	//1.获取url参数
	id := route.GetRouteVariable("id", r)

	//2.读取对应的文章数据
	article, err := article.Get(id)

	//3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//3.1 数据为找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			//3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		//将字符串转化为int64
		id1, _ := strconv.ParseInt(id, 10, 64)
		data := ArticlesFormData{
			Title:  article.Title,
			Body:   article.Body,
			Id:     id1,
			Errors: nil,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logger.LogError(err)

		tmpl.Execute(w, data)
	}
}

// 更新文章
func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	//1.获取URL参数
	id := route.GetRouteVariable("id", r)

	//获取对应的文章数据
	_article, err := article.Get(id)

	//3.如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//3.1数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			//3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		//4. 未出现错误

		//4.1 表单验证
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := validateArticleFormData(title, body)

		if len(errors) == 0 {
			//4.2 表单验证通过，更新数据
			_article.Title = title
			_article.Body = body
			rowsAffected, err := _article.Update()

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
				return
			}

			//更新成功，跳转到文章详情页
			if rowsAffected > 0 {
				http.Redirect(w, r, "/articles/"+id, http.StatusFound)
			} else {
				fmt.Fprint(w, "您没有做任何更改!")
			}
		} else {
			//4.3 表单验证不通过，显示理由
			id1, _ := strconv.ParseInt(id, 10, 64)
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				Id:     id1,
				Errors: errors,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			logger.LogError(err)
			tmpl.Execute(w, data)
		}
	}

}

//delete 删除文章
func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	//1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	//2. 读取对应的文章数据
	_article, err := article.Get(id)

	//3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			//3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		//4. 未出现错误，执行删除操作
		rowsAffected, err := _article.Delete()

		//4.1 发生错误
		if err != nil {
			// 应该是sql 报错了
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			//4.2 未发生错误
			if rowsAffected > 0 {
				//重定向到文章列表页
				indexURL := route.Name2URL("articles.index")
				http.Redirect(w,r,indexURL,http.StatusFound)
			} else {
				//Edga case
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}

}
