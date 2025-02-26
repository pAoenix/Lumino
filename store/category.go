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
func (s *CategoryStore) Register(Category *model.Category) error {
	return s.db.Create(Category).Error
}

// Get -
func (s *CategoryStore) Get(categoryReq *model.CategoryReq) (resp []model.Category, err error) {
	if s.db.Where(categoryReq).Find(&resp).Error != nil {
		return nil, err
	} else {
		return
	}
}
