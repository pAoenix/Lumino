package service

import (
	"Lumino/common/http_error_code"
	"Lumino/model"
	"Lumino/store"
	"mime/multipart"
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
func (s *UserService) Register(userReq *model.RegisterUserReq, fileHeader *multipart.FileHeader) (user model.User, err error) {
	// 1. 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return user, http_error_code.Internal("打开文件失败",
			http_error_code.WithInternal(err))
	}
	defer file.Close()
	// 注册用户
	return s.UserStore.Register(userReq, file)
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
