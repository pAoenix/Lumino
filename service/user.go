package service

import (
	"Lumino/common"
	"Lumino/common/http_error_code"
	"Lumino/common/logger"
	"Lumino/model"
	"Lumino/store"
	"fmt"
	"mime/multipart"
	"sync"
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

// DownloadUserIcons -
func (s *UserService) DownloadUserIcons(userIDs []uint) (users []model.User, err error) {
	users, err = s.UserStore.BatchGetByIDs(userIDs)
	if err != nil {
		return
	}
	const maxConcurrency = 20 // 最大并发数
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup
	errCh := make(chan error, len(users))
	for idx, _ := range users {
		wg.Add(1) // 计数器加1
		go func(i int) {
			sem <- struct{}{} // 获取信号量
			defer func() {
				<-sem // 释放信号量
				wg.Done()
			}()
			if ossUrl, err := s.ossClient.DownloadFile(users[i].IconUrl); err != nil {
				errCh <- fmt.Errorf("处理 %d 失败: %v", i, err)
			} else {
				users[i].IconUrl = ossUrl
			}
		}(idx)
	}
	wg.Wait() // 等待所有goroutine完成
	close(errCh)
	close(sem)
	if len(errCh) != 0 {
		for len(errCh) > 0 {
			logger.Error(<-errCh)
		}
		return nil, http_error_code.Internal("下载用户头像失败")
	}
	return
}
