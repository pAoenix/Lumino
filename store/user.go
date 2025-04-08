package store

import (
	"Lumino/model"
)

// UserStore -
type UserStore struct {
	db *DB
}

// NewUserStore -
func NewUserStore(db *DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

// Register -
func (s *UserStore) Register(user *model.User) error {
	return s.db.Model(model.User{}).Create(user).Error
}

// Modify -
func (s *UserStore) Modify(modifyUserReq *model.ModifyUser) (user model.User, err error) {
	if err = s.db.Model(model.User{}).Where("id = ?", modifyUserReq.ID).Updates(modifyUserReq).Scan(&user).Error; err != nil {
		return user, err
	}
	return
}

// Get -
func (s *UserStore) Get(userReq *model.GetUser) (user model.User, err error) {
	if err = s.db.Model(model.User{}).Where("id = ?", userReq.ID).First(&user).Error; err != nil {
		return user, err
	}
	return
}

// BatchGetByIDs -
func (s *UserStore) BatchGetByIDs(userIDs []int) (users []model.User, err error) {
	if err = s.db.Model(model.User{}).Where("id in ?", userIDs).Find(&users).Error; err != nil {
		return nil, err
	}
	return
}

// Delete -
func (s *UserStore) Delete(userReq *model.DeleteUser) error {
	return s.db.Model(model.User{}).Delete(&model.User{Model: model.Model{ID: userReq.ID}}).Error
}
