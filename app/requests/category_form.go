package requests

import (
	"github.com/thedevsaddam/govalidator"
	"goblog/app/models/category"
)

// ValidateCategoryForm 验证表单，返回 errs 长度等于零即通过
func ValidateCategoryForm(data category.Category) map[string][]string {
	//1. 定制认证规则
	rules := govalidator.MapData{
		"name": []string{"required","min_cn:2","max_cn:8","not_exists:categories,name"},
	}

	//2. 定制错误信息
	messages := govalidator.MapData{
		"name": []string{
			"required:分类名称为必填项",
			"min_cn:分类名称长度需至少 2 个字",
			"max_cn:分类名称长度不能超过8个字",
		},
	}

	//3. 配置初始化
	opts := govalidator.Options{
		Data: &data,
		Rules: rules,
		TagIdentifier: "valid",//模型中的struct 标签标识符
		Messages: messages,
	}

	//4. 开始验证
	return govalidator.New(opts).ValidateStruct()

}
