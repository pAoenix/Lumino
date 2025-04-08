package service

import (
	"Lumino/model"
	"Lumino/store"
)

// CategoryService -
type CategoryService struct {
	CategoryStore *store.CategoryStore
}

// NewCategoryService -
func NewCategoryService(CategoryStore *store.CategoryStore) *CategoryService {
	return &CategoryService{
		CategoryStore: CategoryStore,
	}
}

// Register -
func (s *CategoryService) Register(Category *model.RegisterCategoryReq) (resp model.Category, err error) {
	return s.CategoryStore.Register(Category)
}

// Get -
func (s *CategoryService) Get(categoryReq *model.GetCategoryReq) (resp []model.Category, err error) {
	return s.CategoryStore.Get(categoryReq)
}

// Modify -
func (s *CategoryService) Modify(Category *model.ModifyCategoryReq) (resp model.Category, err error) {
	return s.CategoryStore.Modify(Category)
}

// Delete -
func (s *CategoryService) Delete(Category *model.DeleteCategoryReq) error {
	return s.CategoryStore.Delete(Category)
}
