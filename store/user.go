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
	return s.db.Create(User).Error
}
