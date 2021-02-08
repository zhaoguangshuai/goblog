package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/app/policies"
	"goblog/pkg/auth"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
)

// ArticlesController 文章相关页面
type ArticlesController struct {
	BaseController
}

//Show 文章详情页
func (ac *ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	//1.获取URL参数
	id := route.GetRouteVariable("id", r)

	//2.读取对应的文章数据
	article, err := article.Get(id)

	//3.如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w,err)
	} else {
		//4.读取成功
		view.Render(w,view.D{
			"Article": article,
			"CanModifyArticle":policies.CanModifyArticle(article),
		},"articles.show","articles._article_meta")
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
		view.Render(w,view.D{
			"Articles":articles,
		},"articles.index","articles._article_meta")
	}

}

func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w,view.D{},"articles.create","articles._form_field")
}

// Store 文章创建页面
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

	//1. 初始化数据
	//user_id,_ := strconv.ParseUint(session.Get("uid").(string),10,64)
	user_id := auth.User().ID
	_article := article.Article{
		Title: r.PostFormValue("title"),
		Body: r.PostFormValue("body"),
		UserID: user_id,
	}

	//2. 表单验证
	//errors := requests.ValidateArticleForm(_article)
	// 检查是否有错误
	//if len(errors) == 0 {
		_article.Create()
		if _article.ID > 0 {
			indexURL := route.Name2URL("articles.show", "id",_article.GetStringID())
			http.Redirect(w,r,indexURL,http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章失败，请联系管理员")
		}
	//} else {
	//	view.Render(w,view.D{
	//		"Article":_article,
	//		"Errors": errors,
	//	},"articles.create","articles._form_field")
	//}
}

// Edit 文章更新页面
func (ac *ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	//1.获取url参数
	id := route.GetRouteVariable("id", r)

	//2.读取对应的文章数据
	_article, err := article.Get(id)

	//3. 如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w,err)
	} else {
		//检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w,r)
		} else {
			//4.读取成功，显示编辑文章表单
			view.Render(w,view.D{
				"Article": _article,
				"Errors": nil,
			}, "articles.edit", "articles._form_field")
		}
	}
}

// 更新文章
func (ac *ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	//1.获取URL参数
	id := route.GetRouteVariable("id", r)

	//获取对应的文章数据
	_article, err := article.Get(id)

	//3.如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w,err)
	} else {
		//4. 未出现错误
		//检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w,r)
		} else {
			//4.1 表单验证
			_article.Title = r.PostFormValue("title")
			_article.Body = r.PostFormValue("body")

			//errors := requests.ValidateArticleForm(_article)

			//if len(errors) == 0 {
			//4.2 表单验证通过，更新数据
			rowsAffected, err := _article.Update()

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
				return
			}

			//更新成功，跳转到文章详情页
			if rowsAffected > 0 {
				showURL := route.Name2URL("articles.show","id",id)
				http.Redirect(w, r, showURL, http.StatusFound)
			} else {
				fmt.Fprint(w, "您没有做任何更改!")
			}
			//} else {
			//	//4.3 表单验证不通过，显示理由
			//	view.Render(w,view.D{
			//		"Article": _article,
			//		"Errors": errors,
			//	}, "articles.edit", "articles._form_field")
			//}
		}


	}

}

//delete 删除文章
func (ac *ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	//1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	//2. 读取对应的文章数据
	_article, err := article.Get(id)

	//3. 如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w,err)
	} else {
		//检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w,r)
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

}
