package store

import (
	"Lumino/common"
	"Lumino/common/http_error_code"
	"Lumino/model"
	"errors"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"mime/multipart"
	"strconv"
)

// UserStore -
type UserStore struct {
	db        *DB
	OssClient *common.OssClient
}

// NewUserStore -
func NewUserStore(db *DB, ossClient *common.OssClient) *UserStore {
	return &UserStore{
		db:        db,
		OssClient: ossClient,
	}
}

// Register -
func (s *UserStore) Register(userReq *model.RegisterUserReq, file multipart.File) (user model.User, err error) {
	user.Name = userReq.Name
	// 1.初步注册
	tx := s.db.Begin()
	if err = tx.Model(&model.User{}).Create(&user).Error; err != nil {
		tx.Rollback()
		if IsDuplicateError(err) {
			return user, http_error_code.Conflict("用户名已注册",
				http_error_code.WithInternal(err))
		}
		return user, http_error_code.Internal("注册用户失败",
			http_error_code.WithInternal(err))
	}
	// 2.数据上传
	iconUrl := viper.GetString("oss.profilePhotoDir") + strconv.Itoa(int(user.ID)) + ".jpg"
	if err = s.OssClient.UploadFile(iconUrl, file); err != nil {
		tx.Rollback()
		return user, err
	}
	// 3.更新文件地址
	modifyReq := model.ModifyUserReq{ID: user.ID, IconUrl: iconUrl}
	user.IconUrl = iconUrl
	if err = tx.Model(&model.User{}).Where("id = ?", user.ID).Updates(&modifyReq).Error; err != nil {
		tx.Rollback()
		return user, http_error_code.Internal("注册用户失败",
			http_error_code.WithInternal(err))
	}
	// 4. 提交事务
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return user, http_error_code.Internal("注册用户失败",
			http_error_code.WithInternal(err))
	}
	return
}

// Modify -
func (s *UserStore) Modify(modifyUserReq *model.ModifyUserReq) (user model.User, err error) {
	// 1.判断信息是否正常
	if err = ParamsJudge(s.db, modifyUserReq.DefaultAccountBookID, modifyUserReq.Friend, &modifyUserReq.ID); err != nil {
		return user, err
	}
	if err = s.db.Model(&model.User{}).Where("id = ?", modifyUserReq.ID).Updates(modifyUserReq).First(&user).Error; err != nil {
		if IsDuplicateError(err) {
			return user, http_error_code.Conflict("用户名已注册",
				http_error_code.WithInternal(err))
		}
		return user, err
	}
	return
}

// ModifyProfilePhoto -
func (s *UserStore) ModifyProfilePhoto(modifyUserReq *model.ModifyProfilePhotoReq, file multipart.File) error {
	// 1.判断用户是否存在
	if err := ParamsJudge(s.db, nil, nil, &modifyUserReq.ID); err != nil {
		return err
	}
	// 2.数据上传
	iconUrl := viper.GetString("oss.profilePhotoDir") + strconv.Itoa(int(modifyUserReq.ID)) + ".jpg"
	return s.OssClient.UploadFile(iconUrl, file)
}

// Get -
func (s *UserStore) Get(userReq *model.GetUserReq) (user model.User, err error) {
	if err = s.db.Model(&model.User{}).Where("id = ?", userReq.ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, http_error_code.BadRequest("用户ID不存在")
		}
		return user, http_error_code.Internal("服务内部错误",
			http_error_code.WithInternal(err))
	}
	return
}

// BatchGetByIDs -
func (s *UserStore) BatchGetByIDs(userIDs []int) (users []model.User, err error) {
	if err = s.db.Model(&model.User{}).Where("id in ?", userIDs).Find(&users).Error; err != nil {
		return nil, err
	}
	return
}

// Delete -
func (s *UserStore) Delete(userReq *model.DeleteUserReq) error {
	if err := ParamsJudge(s.db, nil, nil, &userReq.ID); err != nil {
		return err
	}
	return s.db.Model(&model.User{}).Delete(&model.User{Model: model.Model{ID: userReq.ID}}).Error
}
