package model

// Friend -
type Friend struct {
	Inviter uint `json:"inviter" form:"inviter"` // 邀请人
	Invitee uint `json:"invitee" form:"invitee"` // 被邀请人
}
