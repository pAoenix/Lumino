package store

import "Lumino/model"

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
	return s.db.Model(model.User{}).Updates(user).Error
}

// Get -
func (s *UserStore) Get(user *model.User) error {
	return s.db.Model(model.User{}).Where(user).Find(user).Error
}

// BatchGetByIDs -
func (s *UserStore) BatchGetByIDs(userIDs []int) (users []model.User, err error) {
	if err = s.db.Model(model.User{}).Where("id in ?", userIDs).Find(&users).Error; err != nil {
		return nil, err
	}
	return
}
