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
