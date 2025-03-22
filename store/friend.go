package store

import (
	"Lumino/model"
	"errors"
	"fmt"
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
	user := model.User{}
	err := s.db.Model(&model.User{}).Where("? = ANY(friend) and id = ?", friend.Invitee, friend.Inviter).Find(&user).Error
	if err != nil {
		return err
	}
	fmt.Println(user)
	if user.ID != 0 {
		return errors.New("你已存在该好友")
	}
	return s.db.Model(&model.User{}).
		Where("id = ?", friend.Inviter).
		Update("friend", gorm.Expr("array_append(friend, ?)", friend.Invitee)).
		Error
}

// Delete -
func (s *FriendStore) Delete(friend *model.Friend) error {
	return s.db.Model(&model.User{}).
		Where("id = ?", friend.Inviter).
		Update("friend", gorm.Expr("array_remove(friend, ?)", friend.Invitee)).
		Error
}
