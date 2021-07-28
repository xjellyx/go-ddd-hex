package vo

import (
	"time"
)

type BaseRes struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type UserRes struct {
	BaseRes
	UUID     string `json:"uuid,omitempty"`
	Username string `json:"username,omitempty"` // 用户名
	Nickname string `json:"nickname,omitempty"` // 昵称
	IsAdmin  bool   `json:"isAdmin,omitempty"`  // true：是管理员
}
