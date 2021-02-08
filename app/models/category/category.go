package category

import (
	"goblog/app/models"
	"goblog/pkg/model"
	"goblog/pkg/route"
	"goblog/pkg/types"
)

type Category struct {
	models.BaseModel

	Name string `gorm:"type:varchar(255);not null;" valid:"name"`
}

//Link 方法用来生成文章链接
func (c Category) Link() string {
	return route.Name2URL("categories.show","id",c.GetStringID())
}

//Get 通过ID 获取分类
func Get(idstr string) (Category,error) {
	var category Category
	id := types.StringToInt(idstr)
	if err := model.DB.First(&category,id).Error;err != nil {
		return category,err
	}
	return category,nil
}
