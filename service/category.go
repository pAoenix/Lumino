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

// CategoryService -
type CategoryService struct {
	CategoryStore *store.CategoryStore
	ossClient     *common.OssClient
}

// NewCategoryService -
func NewCategoryService(CategoryStore *store.CategoryStore, ossClient *common.OssClient) *CategoryService {
	return &CategoryService{
		CategoryStore: CategoryStore,
		ossClient:     ossClient,
	}
}

// Register -
func (s *CategoryService) Register(Category *model.RegisterCategoryReq, fileHeader *multipart.FileHeader) (resp model.Category, err error) {
	// 1. 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return resp, http_error_code.Internal("打开文件失败",
			http_error_code.WithInternal(err))
	}
	defer file.Close()
	// 2. 注册图标
	return s.CategoryStore.Register(Category, file)
}

// Get -
func (s *CategoryService) Get(categoryReq *model.GetCategoryReq) (resp []model.Category, err error) {
	if resp, err = s.CategoryStore.Get(categoryReq); err != nil {
		return nil, err
	}
	const maxConcurrency = 20 // 最大并发数
	sem := make(chan struct{}, maxConcurrency)
	errCh := make(chan error, len(resp))
	var wg sync.WaitGroup

	defer func() {
		close(errCh)
		close(sem)
	}()
	for idx, _ := range resp {
		wg.Add(1) // 计数器加1
		go func(i int) {
			sem <- struct{}{} // 获取信号量
			defer func() {
				<-sem // 释放信号量
				wg.Done()
			}()
			if ossUrl, err := s.ossClient.DownloadFile(resp[i].IconUrl); err != nil {
				errCh <- fmt.Errorf("处理 %d 失败: %v", i, err)
			} else {
				resp[i].IconUrl = ossUrl
			}
		}(idx)
	}
	wg.Wait() // 等待所有goroutine完成
	if len(errCh) != 0 {
		for len(errCh) > 0 {
			logger.Error(<-errCh)
		}
		return resp, http_error_code.Internal("下载图标异常",
			http_error_code.WithInternal(err))
	}
	return
}

// Modify -
func (s *CategoryService) Modify(Category *model.ModifyCategoryReq) (resp model.Category, err error) {
	return s.CategoryStore.Modify(Category)
}

// ModifyProfilePhoto -
func (s *CategoryService) ModifyProfilePhoto(Category *model.ModifyCategoryIconReq, fileHeader *multipart.FileHeader) error {
	// 1. 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return http_error_code.Internal("打开文件失败",
			http_error_code.WithInternal(err))
	}
	defer file.Close()
	// 2. 修改图标
	return s.CategoryStore.ModifyProfilePhoto(Category, file)
}

// Delete -
func (s *CategoryService) Delete(Category *model.DeleteCategoryReq) error {
	return s.CategoryStore.Delete(Category)
}
