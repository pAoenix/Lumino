package store

import (
	"Lumino/model"
	"fmt"
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
func (s *UserStore) Modify(user *model.User) error {
	return s.db.Model(model.User{}).Where("id = ", user.ID).Updates(user).Error
}

// Get -
func (s *UserStore) Get(user *model.User) (users []model.User, err error) {
	fmt.Println(user)
	if err = s.db.Model(model.User{}).Where(user).Find(&users).Error; err != nil {
		return nil, err
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
