package service

import "Lumino/model"

// UserIconDownloader -
type UserIconDownloader interface {
	DownloadUserIcons(userIDs []uint) ([]model.User, error)
}

// CategoryDownloader -
type CategoryDownloader interface {
	DownloadCategoryIcon(CategoryIDs []uint, categoryInit []model.Category) ([]model.Category, error)
}
