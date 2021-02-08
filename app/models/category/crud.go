package category

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
)

func (Category *Category) Create() (err error) {
	if err = model.DB.Create(&Category).Error;err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}

//all 获取分类数据
func All() ([]Category,error) {
	var categories []Category
	if err := model.DB.Find(&categories).Error;err != nil {
		return categories,err
	}
	return categories,nil
}
