package dao

import (
	"errors"
	"time"

	"github.com/d1vshar/splitgo/model"
	"github.com/google/uuid"
)

func (dao AppDao) GetCategories() (categories []model.Category, err error) {
	err = dao.db.Where("deleted_at IS null").Find(&categories).Error

	return categories, err
}

func (dao AppDao) CreateCategory(Name string) (*model.Category, error) {
	new_category := model.Category{
		Name: Name,
	}

	err := dao.db.Create(&new_category).Error

	return &new_category, err
}

func (dao AppDao) GetCategory(ID uuid.UUID) (category model.Category, err error) {
	err = dao.db.Where("deleted_at IS null").First(&category, ID).Error

	return category, err
}

func (dao AppDao) UpdateCategory(ID uuid.UUID, Name string) (*model.Category, error) {
	category, err := dao.GetCategory(ID)
	if err != nil {
		return nil, err
	}

	category.Name = Name

	return &category, dao.db.Save(&category).Error
}

// soft delete
func (dao AppDao) DeleteCategory(ID uuid.UUID) error {
	res, err := dao.GetCategory(ID)

	if err != nil {
		return errors.New("category not found")
	}

	now := time.Now()
	res.Meta.DeletedAt = &now

	return dao.db.Save(&res).Error
}
