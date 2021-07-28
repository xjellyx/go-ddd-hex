package entity

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserUUID string `gorm:"not null;type:varchar(36)"`
	Title    string `gorm:"type:varchar(64)"` // 文章标题
	Content  string // 文章内容
}

func (p *Post) SetUserUUID(userID string) *Post {
	p.UserUUID = userID
	return p
}

func (p *Post) QueryCond() (res map[string]interface{}) {
	res = map[string]interface{}{}
	switch {
	case len(p.UserUUID) > 0:
		res["user_uuid"] = p.UserUUID
	case len(p.Title) > 0:
		res["title"] = p.Title

	}
	return
}
