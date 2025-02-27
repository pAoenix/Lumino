package model

import (
	"gorm.io/gorm"
	"time"
)

// Transaction 交易记录
type Transaction struct {
	gorm.Model
	Type          int       // 类型:收入/支出
	Amount        float32   //交易数额
	Date          time.Time //日期
	UserId        uint      // 账户id
	CategoryId    uint      //关联消费场景分类ID
	Description   string    //注释
	AccountBookID int       // 对应的账本id
}

type TransactionReq struct {
	UserId        uint // 账户id
	AccountBookID int  // 对应的账本id
}
