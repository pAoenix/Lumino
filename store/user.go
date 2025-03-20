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
func (s *UserStore) Register(User *model.User) error {
	return s.db.Model(model.User{}).Create(User).Error
}

// Modify -
func (s *UserStore) Modify(User *model.User) error {
	return s.db.Model(model.User{}).Updates(User).Error
}

// Get -
func (s *UserStore) Get(User *model.User) error {
	return s.db.Model(model.User{}).Where(User).Find(User).Error
}

// BatchGetByIDs -
func (s *UserStore) BatchGetByIDs(userIDs []int) (users []model.User, err error) {
	if err = s.db.Model(model.User{}).Where("id in ?", userIDs).Find(&users).Error; err != nil {
		return nil, err
	}
	return
}
