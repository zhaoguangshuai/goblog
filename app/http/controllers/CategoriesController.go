package controllers

import (
	"fmt"
	"goblog/app/models/category"
	"goblog/app/requests"
	"goblog/pkg/view"
	"net/http"
)
//文章分类控制器
type CategoriesController struct {
	BaseController
}
// 文章分类创建页面
func (*CategoriesController) Create(w http.ResponseWriter,r *http.Request)  {
	view.Render(w,view.D{},"categories.create")
}

//保存文章分类
func (*CategoriesController) Store(w http.ResponseWriter,r *http.Request)  {
	//1. 初始化数据
	_category := category.Category{
		Name: r.PostFormValue("name"),
	}

	//表单验证
	errors := requests.ValidateCategoryForm(_category)

	//3. 检测错误
	if len(errors) == 0 {
		//创建文章分类
		_category.Create()
		if _category.ID > 0 {
			fmt.Fprint(w, "创建成功！")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"创建文章分类失败，请联系管理员")
		}
	} else {
		view.Render(w,view.D{
			"Category": _category,
			"Errors":errors,
		},"categories.create")
	}
}
