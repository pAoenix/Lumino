package service

import (
	"Lumino/common"
	"Lumino/model"
	"Lumino/store"
	"gorm.io/gorm"
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
	user := &model.User{Model: gorm.Model{ID: accountBookReq.UserId}}
	err = s.UserStore.Get(user)
	if err != nil {
		return
	}
	resp.DefaultAccountBookID = user.DefaultAccountBookID
	// 计算涉及的用户信息
	var userIDs []int
	for _, abl := range accountBookList {
		for _, userID := range abl.UserId {
			if !common.ContainsInt(userIDs, userID) {
				userIDs = append(userIDs, userID)
			}
		}
	}
	users, err := s.UserStore.BatchGetByIDs(userIDs)
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
