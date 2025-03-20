package store

import (
	"Lumino/model"
	"gorm.io/gorm"
)

// FriendStore -
type FriendStore struct {
	db *DB
}

// NewFriendStore -
func NewFriendStore(db *DB) *FriendStore {
	return &FriendStore{
		db: db,
	}
}

// Invite -
func (s *FriendStore) Invite(friend *model.Friend) error {
	return s.db.Model(&model.User{}).
		Update("friend", gorm.Expr("array_append(friend, ?)", friend.Invitee)).
		Where("id = ?", friend.Inviter).
		Error
}

// Delete -
func (s *FriendStore) Delete(friend *model.Friend) error {
	return s.db.Model(&model.User{}).
		Update("friend", gorm.Expr("array_remove(friend, ?)", friend.Invitee)).
		Where("id = ?", friend.Inviter).
		Error
}
