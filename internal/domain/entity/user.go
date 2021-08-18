package entity

import (
	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
)

type User struct {
	gorm.Model
	UUID     string      `gorm:"uniqueIndex;not null;type:varchar(36)"`
	Username string      `gorm:"uniqueIndex;not null;type:varchar(36)"` // 用户名
	Password null.String `gorm:"type:varchar(256)"`                     // 密码
	Nickname null.String `gorm:"type:varchar(36)"`                      // 昵称
	IsAdmin  null.Bool   `gorm:"default: false"`                        // true：是管理员
}

func NewUser(username string) *User {
	return &User{
		UUID:     uuid.NewV4().String(),
		Username: username,
	}
}

func (u *User) SetID(id string) *User {
	_id, _ := strconv.Atoi(id)
	u.ID = uint(_id)
	return u
}

func (u *User) SetUUID(uuid string) *User {
	u.UUID = uuid
	return u
}

func (u *User) SetUsername(username string) *User {
	u.Username = username
	return u
}

func (u *User) SetNickname(val string) *User {
	if len(val) > 0 {
		u.Nickname.SetValid(val)
	}
	return u
}

func (u *User) SetPassword(val string) *User {
	if len(val) > 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte(val), bcrypt.DefaultCost)
		u.Password.SetValid(string(hash))
	}
	return u
}

func (u *User) SetIsAdmin(val *bool) *User {
	if val != nil {
		u.IsAdmin.SetValid(*val)
	}
	return u
}

func (u *User) QueryCond() (res map[string]interface{}) {
	var (
		data = make(map[string]interface{})
	)
	switch {
	// 优先检测唯一字段
	case u.ID > 0:
		data["id"] = u.ID
		return data
	case len(u.UUID) > 0:
		data["uuid"] = u.UUID
		return data
	case len(u.Username) > 0:
		data["username"] = u.Username
		return data
	case u.Nickname.Ptr() != nil:
		data["nickname"] = u.Nickname.String
	case u.IsAdmin.Ptr() != nil:
		data["is_admin"] = u.IsAdmin.Bool
	}
	return data
}
