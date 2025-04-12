package service

import (
	"Lumino/common"
	"Lumino/common/logger"
	"Lumino/model"
	"Lumino/store"
	"errors"
	"fmt"
	"sync"
)

// AccountBookService -
type AccountBookService struct {
	AccountBookStore *store.AccountBookStore
	UserStore        *store.UserStore
	ossClient        *common.OssClient
}

// NewAccountBookService -
func NewAccountBookService(accountBookStore *store.AccountBookStore, userStore *store.UserStore, ossClient *common.OssClient) *AccountBookService {
	return &AccountBookService{
		AccountBookStore: accountBookStore,
		UserStore:        userStore,
		ossClient:        ossClient,
	}
}

// Register -
func (s *AccountBookService) Register(accountBook *model.RegisterAccountBookReq) error {
	if !common.ContainsInt(common.ConvertArrayToIntSlice(accountBook.UserIDs), int(accountBook.CreatorID)) {
		accountBook.UserIDs = append(accountBook.UserIDs, int32(accountBook.CreatorID))
	}
	return s.AccountBookStore.Register(accountBook)
}

// Modify -
func (s *AccountBookService) Modify(accountBook *model.ModifyAccountBookReq) error {
	return s.AccountBookStore.Modify(accountBook)
}

// Get -
func (s *AccountBookService) Get(accountBookReq *model.GetAccountBookReq) (resp model.AccountBookResp, err error) {
	// 账本汇总
	accountBookList, err := s.AccountBookStore.Get(accountBookReq)
	if err != nil {
		return
	}
	resp.AccountBooks = accountBookList

	// 计算默认账本
	userReq := &model.GetUserReq{ID: accountBookReq.UserID}
	user, err := s.UserStore.Get(userReq)
	if err != nil {
		return
	}
	resp.DefaultAccountBookID = user.DefaultAccountBookID
	// 计算涉及的用户信息
	var userIDs []int
	for _, abl := range accountBookList {
		for _, userID := range abl.UserIDs {
			userIDint := int(userID)
			if !common.ContainsInt(userIDs, userIDint) {
				userIDs = append(userIDs, userIDint)
			}
		}
	}
	users, err := s.UserStore.BatchGetByIDs(userIDs)
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
		return resp, errors.New("下载用户头像失败")
	}
	resp.Users = users
	return
}

// Delete -
func (s *AccountBookService) Delete(accountBook *model.DeleteAccountBookReq) error {
	return s.AccountBookStore.Delete(accountBook)
}

// Merge -
func (s *AccountBookService) Merge(mergeAccountBookReq *model.MergeAccountBookReq) error {
	return s.AccountBookStore.Merge(mergeAccountBookReq)
}
