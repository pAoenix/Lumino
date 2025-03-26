package service

import (
	"Lumino/model"
	"Lumino/store"
)

// UserService -
type UserService struct {
	UserStore *store.UserStore
}

// NewUserService -
func NewUserService(UserStore *store.UserStore) *UserService {
	return &UserService{
		UserStore: UserStore,
	}
}

// Register -
func (s *UserService) Register(user *model.User) error {
	return s.UserStore.Register(user)
}

// Modify -
func (s *UserService) Modify(user *model.User) error {
	return s.UserStore.Modify(user)
}

// Get -
func (s *UserService) Get(user *model.User) (users []model.User, err error) {
	return s.UserStore.Get(user)
}
