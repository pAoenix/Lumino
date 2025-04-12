package service

import (
	"Lumino/common"
	"Lumino/common/http_error_code"
	"Lumino/model"
	"Lumino/store"
	"mime/multipart"
)

// UserService -
type UserService struct {
	UserStore *store.UserStore
	ossClient *common.OssClient
}

// NewUserService -
func NewUserService(UserStore *store.UserStore, ossClient *common.OssClient) *UserService {
	return &UserService{
		UserStore: UserStore,
		ossClient: ossClient,
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

// ModifyProfilePhoto -
func (s *UserService) ModifyProfilePhoto(userReq *model.ModifyProfilePhotoReq, fileHeader *multipart.FileHeader) error {
	// 1. 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return http_error_code.Internal("打开文件失败",
			http_error_code.WithInternal(err))
	}
	defer file.Close()
	// 2. 修改图标
	return s.UserStore.ModifyProfilePhoto(userReq, file)
}

// Modify -
func (s *UserService) Modify(modifyUserReq *model.ModifyUserReq) (user model.User, err error) {
	return s.UserStore.Modify(modifyUserReq)
}

// Get -
func (s *UserService) Get(userReq *model.GetUserReq) (user model.User, err error) {
	if user, err = s.UserStore.Get(userReq); err != nil {
		return
	}
	if ossUrl, err := s.ossClient.DownloadFile(user.IconUrl); err != nil {
		return user, err
	} else {
		user.IconUrl = ossUrl
	}
	return
}

// Delete -
func (s *UserService) Delete(userReq *model.DeleteUserReq) error {
	return s.UserStore.Delete(userReq)
}
