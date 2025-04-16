package service

import (
	"Lumino/model"
	"Lumino/store"
)

// FriendService -
type FriendService struct {
	FriendStore *store.FriendStore
}

// NewFriendService -
func NewFriendService(FriendStore *store.FriendStore) *FriendService {
	return &FriendService{
		FriendStore: FriendStore,
	}
}

// Invite -
func (s *FriendService) Invite(friend *model.Friend) (user model.User, err error) {
	return s.FriendStore.Invite(friend)
}

// Delete -
func (s *FriendService) Delete(friend *model.Friend) (user model.User, err error) {
	return s.FriendStore.Delete(friend)
}
