package store

import "Lumino/model"

// CategoryStore -
type CategoryStore struct {
	db *DB
}

// NewCategoryStore -
func NewCategoryStore(db *DB) *CategoryStore {
	return &CategoryStore{
		db: db,
	}
}

// Register -
func (s *CategoryStore) Register(category *model.RegisterCategoryReq) (resp model.Category, err error) {
	resp.Name = category.Name
	resp.UserID = category.UserID
	user := model.User{}
	// 判断用户是否存在
	if err = s.db.Model(&model.User{}).Where("id = ?", category.UserID).First(&user).Error; err != nil {
		return
	}
	err = s.db.Model(&model.Category{}).Create(&resp).Error
	return
}

// Get -
func (s *CategoryStore) Get(categoryReq *model.GetCategoryReq) (resp []model.Category, err error) {
	err = s.db.Model(&model.Category{}).Where(categoryReq).Find(&resp).Error
	return
}

// Modify -
func (s *CategoryStore) Modify(category *model.ModifyCategoryReq) (resp model.Category, err error) {
	err = s.db.Model(&model.Category{}).Where("id = ?", category.ID).Updates(category).Find(&resp).Error
	return
}

// Delete -
func (s *CategoryStore) Delete(category *model.DeleteCategoryReq) error {
	return s.db.Model(&model.Category{}).Delete(&model.Category{Model: model.Model{ID: category.ID}}).Error
}
