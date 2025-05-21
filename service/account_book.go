package service

import (
	"Lumino/common"
	"Lumino/model"
	"Lumino/store"
	"math"
	"sort"
)

// AccountBookService -
type AccountBookService struct {
	AccountBookStore *store.AccountBookStore
	UserStore        *store.UserStore
	TransactionStore *store.TransactionStore
	UserDownloader   UserIconDownloader
}

// NewAccountBookService -
func NewAccountBookService(accountBookStore *store.AccountBookStore, userStore *store.UserStore, transactionStore *store.TransactionStore, userService UserIconDownloader) *AccountBookService {
	return &AccountBookService{
		AccountBookStore: accountBookStore,
		UserStore:        userStore,
		TransactionStore: transactionStore,
		UserDownloader:   userService,
	}
}

// Balance 用于分账系统临时使用
type Balance struct {
	UserID uint
	Amount float64
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
	if users, err := s.UserDownloader.DownloadUserIcons(userIDs); err != nil {
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

// AA -
func (s *AccountBookService) AA(accountBook *model.AAAccountBookReq) (transfers []model.AAResult, err error) {
	// 获取账本相关的用户记录
	req := model.GetTransactionReq{
		AccountBookID: &accountBook.ID,
		UserID:        &accountBook.UserID,
	}
	transactions, err := s.TransactionStore.Get(&req)

	if err != nil {
		return
	}
	// 开始分账
	balanceMap := make(map[uint]float64)
	for _, transaction := range transactions {
		// 计算AA后每人应承担的金额
		splitAmount := transaction.Amount / float64(len(transaction.RelatedUserIDs))

		// 付款人支付了全额，所以先加回全额
		balanceMap[transaction.PayUserID] += transaction.Amount

		// 每个相关用户(包括付款人自己)需要承担splitAmount
		for _, userID := range transaction.RelatedUserIDs {
			balanceMap[uint(userID)] -= splitAmount
		}
	}
	// 2. 将余额转换为列表并分离债权人和债务人
	var balances []Balance
	for userID, amount := range balanceMap {
		// 忽略余额为0的用户
		if math.Abs(amount) > model.Epsilon { // 使用小的epsilon处理浮点精度
			balances = append(balances, Balance{UserID: userID, Amount: amount})
		}
	}
	// 3. 分离债权人和债务人
	var creditors, debtors []Balance
	for _, b := range balances {
		if b.Amount > 0 {
			creditors = append(creditors, b)
		} else {
			debtors = append(debtors, Balance{UserID: b.UserID, Amount: -b.Amount})
		}
	}
	// 4. 排序：金额大的优先处理
	sort.Slice(creditors, func(i, j int) bool {
		return creditors[i].Amount > creditors[j].Amount
	})
	sort.Slice(debtors, func(i, j int) bool {
		return debtors[i].Amount > debtors[j].Amount
	})
	// 5. 贪心算法匹配债权人和债务人
	cIdx, dIdx := 0, 0

	for cIdx < len(creditors) && dIdx < len(debtors) {
		creditor := creditors[cIdx]
		debtor := debtors[dIdx]

		// 计算可以转账的金额
		amount := math.Min(creditor.Amount, debtor.Amount)

		// 创建转账记录
		transfers = append(transfers, model.AAResult{
			DebtorID:   debtor.UserID,
			CreditorID: creditor.UserID,
			Amount:     amount,
		})

		// 更新余额
		creditors[cIdx].Amount -= amount
		debtors[dIdx].Amount -= amount

		// 移动指针
		if creditors[cIdx].Amount < model.Epsilon {
			cIdx++
		}
		if debtors[dIdx].Amount < model.Epsilon {
			dIdx++
		}
	}
	if transfers == nil {
		transfers = []model.AAResult{}
	}
	return
}
