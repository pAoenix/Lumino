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
func (s *UserService) Modify(modifyUserReq *model.ModifyUser) (user model.User, err error) {
	return s.UserStore.Modify(modifyUserReq)
}

// Get -
func (s *UserService) Get(userReq *model.GetUser) (user model.User, err error) {
	return s.UserStore.Get(userReq)
}

// Delete -
func (s *UserService) Delete(userReq *model.DeleteUser) error {
	return s.UserStore.Delete(userReq)
}
