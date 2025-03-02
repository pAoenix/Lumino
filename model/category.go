package model

import "gorm.io/gorm"

// Category 记账类别
type Category struct {
	gorm.Model
	Name          string //分类名称的中文示意
	UserID        uint   // 用户id
	IconUrl       string // 类别图标的本地地址
	IconRemoteUrl string // 类别图标的对象存储下载地址
}

// CategoryReq -
type CategoryReq struct {
	UserID uint // 用户id
}
