package model

// Friend -
type Friend struct {
	Inviter int `json:"inviter" form:"inviter"` // 邀请人
	Invitee int `json:"invitee" form:"invitee"` // 被邀请人
}
