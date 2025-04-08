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
func (s *UserService) Register(userReq *model.RegisterUserReq) (user model.User, err error) {
	return s.UserStore.Register(userReq)
}

// Modify -
func (s *UserService) Modify(modifyUserReq *model.ModifyUserReq) (user model.User, err error) {
	return s.UserStore.Modify(modifyUserReq)
}

// Get -
func (s *UserService) Get(userReq *model.GetUserReq) (user model.User, err error) {
	return s.UserStore.Get(userReq)
}

// Delete -
func (s *UserService) Delete(userReq *model.DeleteUserReq) error {
	return s.UserStore.Delete(userReq)
}
