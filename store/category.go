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
func (s *CategoryStore) Register(category *model.Category) error {
	return s.db.Model(&model.Category{}).Create(category).Error
}

// Get -
func (s *CategoryStore) Get(categoryReq *model.CategoryReq) (resp []model.Category, err error) {
	if s.db.Model(&model.Category{}).Where(categoryReq).Find(&resp).Error != nil {
		return nil, err
	} else {
		return
	}
}

// Modify -
func (s *CategoryStore) Modify(category *model.Category) error {
	return s.db.Model(&model.Category{}).Where("id = ?", category.Model.ID).Updates(category).Error
}
