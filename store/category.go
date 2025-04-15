package store

import (
	"Lumino/common"
	"Lumino/common/http_error_code"
	"Lumino/model"
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
	"mime/multipart"
	"strconv"
)

// CategoryStore -
type CategoryStore struct {
	db        *DB
	ossClient *common.OssClient
}

// NewCategoryStore -
func NewCategoryStore(db *DB, ossClient *common.OssClient) *CategoryStore {
	return &CategoryStore{
		db:        db,
		ossClient: ossClient,
	}
}

// Register -
func (s *CategoryStore) Register(categoryReq *model.RegisterCategoryReq, file multipart.File) (category model.Category, err error) {
	if err = ParamsJudge(s.db, nil, nil, &categoryReq.UserID, nil, nil); err != nil {
		return category, err
	}
	if err = copier.Copy(&category, &categoryReq); err != nil {
		return category, http_error_code.Internal("服务内异常",
			http_error_code.WithInternal(err))
	}
	tx := s.db.Begin()
	if err = tx.Model(&model.Category{}).Create(&category).Error; err != nil {
		tx.Rollback()
		if IsDuplicateError(err) {
			return category, http_error_code.Conflict("已注册同名图标",
				http_error_code.WithInternal(err))
		}
		return category, http_error_code.Internal("注册图标失败",
			http_error_code.WithInternal(err))
	}
	// 2.数据上传
	iconUrl := viper.GetString("oss.categoryDir") + strconv.Itoa(int(category.ID)) + ".jpg"
	if err = s.ossClient.UploadFile(iconUrl, file); err != nil {
		tx.Rollback()
		return category, err
	}
	// 3.更新文件地址
	modifyReq := model.ModifyCategoryReq{ID: category.ID, IconUrl: iconUrl}
	category.IconUrl = iconUrl
	if err = tx.Model(&model.Category{}).Where("id = ?", category.ID).Updates(&modifyReq).Error; err != nil {
		tx.Rollback()
		return category, http_error_code.Internal("注册用户失败",
			http_error_code.WithInternal(err))
	}
	// 4. 提交事务
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return category, http_error_code.Internal("注册用户失败",
			http_error_code.WithInternal(err))
	}
	return
}

// Get -
func (s *CategoryStore) Get(categoryReq *model.GetCategoryReq) (resp []model.Category, err error) {
	if err = ParamsJudge(s.db, nil, nil, categoryReq.UserID, categoryReq.ID, nil); err != nil {
		return resp, err
	}
	err = s.db.Model(&model.Category{}).Where(categoryReq).Find(&resp).Error
	return
}

// Modify -
func (s *CategoryStore) Modify(categoryReq *model.ModifyCategoryReq) (resp model.Category, err error) {
	if err = ParamsJudge(s.db, nil, nil, categoryReq.UserID, &categoryReq.ID, nil); err != nil {
		return resp, err
	}
	category := model.Category{}
	if err = copier.Copy(&category, &categoryReq); err != nil {
		return resp, http_error_code.Internal("服务内异常",
			http_error_code.WithInternal(err))
	}
	err = s.db.Model(&model.Category{}).Where("id = ?", categoryReq.ID).Updates(category).Find(&resp).Error
	return
}

// ModifyProfilePhoto -
func (s *CategoryStore) ModifyProfilePhoto(categoryReq *model.ModifyCategoryIconReq, file multipart.File) error {
	if err := ParamsJudge(s.db, nil, nil, nil, &categoryReq.ID, nil); err != nil {
		return err
	}
	// 2.数据上传
	iconUrl := viper.GetString("oss.categoryDir") + strconv.Itoa(int(categoryReq.ID)) + ".jpg"
	return s.ossClient.UploadFile(iconUrl, file)
}

// Delete -
func (s *CategoryStore) Delete(category *model.DeleteCategoryReq) error {
	if err := ParamsJudge(s.db, nil, nil, nil, &category.ID, nil); err != nil {
		return err
	}
	return s.db.Model(&model.Category{}).Delete(&model.Category{Model: model.Model{ID: category.ID}}).Error
}
