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
func (s *UserService) Register(User *model.User) error {
	return s.UserStore.Register(User)
}

// Modify -
func (s *UserService) Modify(User *model.User) error {
	return s.UserStore.Modify(User)
}

// Get -
func (s *UserService) Get(User *model.User) error {
	return s.UserStore.Get(User)
}
