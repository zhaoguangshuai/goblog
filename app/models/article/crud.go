package article

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/types"
	"net/http"
)
// Get 通过 ID 获取文章
func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToInt(idstr)
	if err := model.DB.Preload("Category").Preload("User").First(&article, id).Error; err != nil {
		return article, err
	}
	return article, nil
}
// GetAll 获取全部文章
func GetAll() ([]Article, error) {
	var articles []Article
	if err := model.DB.Preload("User").Preload("Category").Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

// Create 创建文章，通过 article.ID 来判断是否创建成功
func (article *Article) Create() (err error) {
	if err = model.DB.Create(&article).Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}

//Update 修改文章信息
func (article *Article) Update() (rowsAffected int64, err error) {
	result := model.DB.Save(&article)
	if err = result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return result.RowsAffected, nil
}

// delete 删除文章
func (article *Article) Delete() (rowsAffected int64, err error) {
	result := model.DB.Delete(&article)
	if err = result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return result.RowsAffected, nil
}

//GetByUserID 获取全部文章
func GetByUserID(uid string) ([]Article,error) {
	var articles []Article
	if err := model.DB.Where("user_id = ?",uid).Preload("User").Find(&articles).Error; err != nil {
		return articles,err
	}
	return articles,nil

}

//GetByCategoryID 获取分类相关的文章
func GetByCategoryID(cid string,r *http.Request) ([]Article,error) {
	var articles []Article
	if err := model.DB.Preload("User").Where("category_id = ?", cid).Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles,nil
}


