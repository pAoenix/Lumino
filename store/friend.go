package store

import (
	"Lumino/common/http_error_code"
	"Lumino/model"
	"errors"
	"github.com/lib/pq"
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
	if err := ParamsJudge(s.db, nil, &pq.Int32Array{int32(friend.Invitee)}, &friend.Inviter); err != nil {
		return err
	}
	user := model.User{}
	if err := s.db.Model(&model.User{}).Where("? = ANY(friend) and id = ?", friend.Invitee, friend.Inviter).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return http_error_code.Internal("服务内部错误",
				http_error_code.WithInternal(err))
		}
	} else {
		return http_error_code.Conflict("你已存在该好友",
			http_error_code.WithInternal(err))
	}
	return s.db.Model(&model.User{}).
		Where("id = ?", friend.Inviter).
		Update("friend", gorm.Expr("array_append(friend, ?)", friend.Invitee)).
		Error
}

// Delete -
func (s *FriendStore) Delete(friend *model.Friend) error {
	if err := ParamsJudge(s.db, nil, &pq.Int32Array{int32(friend.Invitee)}, &friend.Inviter); err != nil {
		return err
	}
	user := model.User{}
	if err := s.db.Model(&model.User{}).Where("? = ANY(friend) and id = ?", friend.Invitee, friend.Inviter).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http_error_code.BadRequest("对方不是你好友",
				http_error_code.WithInternal(err))
		}
		return http_error_code.Internal("服务内部错误",
			http_error_code.WithInternal(err))
	}
	return s.db.Model(&model.User{}).
		Where("id = ?", friend.Inviter).
		Update("friend", gorm.Expr("array_remove(friend, ?)", friend.Invitee)).
		Error
}
