package model

// Category 记账类别
type Category struct {
	Model
	Name            string `json:"name" form:"name"`                           //分类名称的中文示意
	UserID          uint   `json:"user_id" form:"user_id"`                     // 用户id
	Icon            string `json:"icon" form:"icon"`                           // 类别图标的本地地址
	IconDownloadUrl string `json:"icon_download_url" form:"icon_download_url"` // 类别图标的对象存储下载地址
}

// CategoryReq -
type CategoryReq struct {
	UserID uint `json:"user_id" form:"user_id"` // 用户id
}
