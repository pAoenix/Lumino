package model

const CategoryTableName = "categories"

// Category 记账类别
type Category struct {
	Model
	Name    string `json:"name" form:"name" gorm:"uniqueIndex:idx_name_userid"`       //分类名称的中文示意
	UserID  uint   `json:"user_id" form:"user_id" gorm:"uniqueIndex:idx_name_userid"` // 用户id
	IconUrl string `json:"icon_url" form:"icon_url"`                                  // 类别图标的本地地址
}

// GetCategoryReq -
type GetCategoryReq struct {
	ID     uint `json:"id" form:"id"`                              // ID 图标id
	UserID uint `json:"user_id" form:"user_id" binding:"required"` // 用户id
}

// DeleteCategoryReq -
type DeleteCategoryReq struct {
	ID uint `json:"id" form:"id" binding:"required"` // ID 图标id
}

// ModifyCategoryReq -
type ModifyCategoryReq struct {
	ID      uint   `json:"id" form:"id" binding:"required"`               // ID 图标id
	Name    string `json:"name" form:"name"`                              //分类名称的中文示意
	UserID  *uint  `json:"user_id" form:"user_id"`                        // 用户id
	IconUrl string `json:"icon_url" form:"icon_url" swaggerignore:"true"` // 类别图标的本地地址
}

// ModifyCategoryIconReq -
type ModifyCategoryIconReq struct {
	ID uint `json:"id" form:"id" binding:"required"` // ID 图标id
}

// RegisterCategoryReq -
type RegisterCategoryReq struct {
	Name   string `json:"name" form:"name" binding:"required"`       //分类名称的中文示意
	UserID uint   `json:"user_id" form:"user_id" binding:"required"` // 用户id
}
