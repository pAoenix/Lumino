package model

// Friend -
type Friend struct {
	Inviter uint `json:"inviter" form:"inviter" binding:"required"` // 邀请人
	Invitee uint `json:"invitee" form:"invitee" binding:"required"` // 被邀请人
}
