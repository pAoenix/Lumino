package service

import (
	"Lumino/common"
	"Lumino/model"
	"Lumino/store"
)

// AccountBookService -
type AccountBookService struct {
	AccountBookStore *store.AccountBookStore
	UserStore        *store.UserStore
	userDownloader   UserIconDownloader
}

// NewAccountBookService -
func NewAccountBookService(accountBookStore *store.AccountBookStore, userStore *store.UserStore, userService UserIconDownloader) *AccountBookService {
	return &AccountBookService{
		AccountBookStore: accountBookStore,
		UserStore:        userStore,
		userDownloader:   userService,
	}
}

// Register -
func (s *AccountBookService) Register(accountBookreq *model.RegisterAccountBookReq) (accountBook model.AccountBook, err error) {
	if !common.ContainsInt(common.ConvertArrayToIntSlice(accountBookreq.UserIDs), int(accountBookreq.CreatorID)) {
		accountBookreq.UserIDs = append(accountBookreq.UserIDs, int32(accountBookreq.CreatorID))
	}
	return s.AccountBookStore.Register(accountBookreq)
}

// Modify -
func (s *AccountBookService) Modify(accountBookReq *model.ModifyAccountBookReq) (accountBook model.AccountBook, err error) {
	return s.AccountBookStore.Modify(accountBookReq)
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
	userReq := &model.GetUserReq{ID: &accountBookReq.CreatorID}
	user, err := s.UserStore.Get(userReq)
	if err != nil {
		return
	}
	resp.DefaultAccountBookID = user.DefaultAccountBookID
	// 计算涉及的用户信息
	var userIDs []uint
	for _, abl := range accountBookList {
		for _, userID := range abl.UserIDs {
			if !common.ContainsUint(userIDs, uint(userID)) {
				userIDs = append(userIDs, uint(userID))
			}
		}
	}
	if users, err := s.userDownloader.DownloadUserIcons(userIDs); err != nil {
		return resp, err
	} else {
		resp.Users = users
	}
	return
}

// Delete -
func (s *AccountBookService) Delete(accountBook *model.DeleteAccountBookReq) error {
	return s.AccountBookStore.Delete(accountBook)
}

// Merge -
func (s *AccountBookService) Merge(mergeAccountBookReq *model.MergeAccountBookReq) (resp model.AccountBookResp, err error) {
	if err = s.AccountBookStore.Merge(mergeAccountBookReq); err != nil {
		return resp, err
	}
	return s.Get(&model.GetAccountBookReq{
		ID:        &mergeAccountBookReq.MergeAccountBookID,
		CreatorID: mergeAccountBookReq.CreatorID,
	})
}
