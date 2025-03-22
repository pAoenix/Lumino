package service

import (
	"Lumino/common"
	"Lumino/model"
	"Lumino/store"
	"errors"
)

// AccountBookService -
type AccountBookService struct {
	AccountBookStore *store.AccountBookStore
	UserStore        *store.UserStore
}

// NewAccountBookService -
func NewAccountBookService(accountBookStore *store.AccountBookStore, userStore *store.UserStore) *AccountBookService {
	return &AccountBookService{
		AccountBookStore: accountBookStore,
		UserStore:        userStore,
	}
}

// Register -
func (s *AccountBookService) Register(accountBook *model.AccountBook) error {
	return s.AccountBookStore.Register(accountBook)
}

// Modify -
func (s *AccountBookService) Modify(accountBook *model.AccountBook) error {
	return s.AccountBookStore.Modify(accountBook)
}

// Get -
func (s *AccountBookService) Get(accountBookReq *model.AccountBookReq) (resp model.AccountBookResp, err error) {
	// 账本汇总
	accountBookList, err := s.AccountBookStore.Get(accountBookReq)
	if err != nil {
		return
	}
	resp.AccountBooks = accountBookList

	// 计算默认账本
	user := &model.User{Model: model.Model{ID: accountBookReq.UserID}}
	users, err := s.UserStore.Get(user)
	if err != nil {
		return
	}
	if len(users) != 1 {
		return resp, errors.New("user-id is error")
	}
	resp.DefaultAccountBookID = users[0].DefaultAccountBookID
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
	users, err = s.UserStore.BatchGetByIDs(userIDs)
	if err != nil {
		return
	}
	resp.Users = users
	return
}

// Delete -
func (s *AccountBookService) Delete(accountBook *model.AccountBook) error {
	return s.AccountBookStore.Delete(accountBook)
}

// Merge -
func (s *AccountBookService) Merge(mergeAccountBookReq *model.MergeAccountBookReq) error {
	return s.AccountBookStore.Merge(mergeAccountBookReq)
}
