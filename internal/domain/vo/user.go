package vo

import (
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	uuid "github.com/satori/go.uuid"
	"strconv"
)

type UserVO struct {
	BaseVO
	UUID     string `json:"uuid,omitempty"`
	Username string `json:"username,omitempty"` // 用户名
	Nickname string `json:"nickname,omitempty"` // 昵称
	Phone    string `json:"phone,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	IsAdmin  bool   `json:"is_admin,omitempty"` // true：是管理员
	//
	Password string `json:"-"`
}

type RegisterForm struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUnique struct {
	Username string `json:"username,omitempty"` // 用户名
	Phone    string `json:"phone,omitempty"`
	UUID     string `json:"uuid,omitempty"`
	ID       string `json:"id,omitempty"`
}

type LoginForm struct {
	Username  string `json:"username,omitempty"` // 用户名
	Phone     string `json:"phone,omitempty"`
	Password  string `json:"password,omitempty"`
	CaptchaId string `form:"captchaId" json:"captchaId" `
	Digits    string `form:"digits" json:"digits"`
}

type UserVOForm struct {
	Username string `json:"username,omitempty"` // 用户名
	Nickname string `json:"nickname,omitempty"` // 昵称
	IsAdmin  bool   `json:"is_admin,omitempty"` // true：是管理员
	Password string `json:"password,omitempty"`
}

func UserEntity2VO(in *entity.User) *UserVO {
	res := new(UserVO)
	res.ID = strconv.Itoa(int(in.ID))
	res.Username = in.Username
	res.Nickname = in.Nickname.String
	res.IsAdmin = in.IsAdmin.Bool
	res.CreatedAt = in.CreatedAt
	res.UpdatedAt = in.UpdatedAt
	res.Password = in.Password.String
	res.Phone = in.Phone
	res.Avatar = in.Avatar.String
	return res
}

func UserVOForm2Entity(in *UserVOForm) *entity.User {
	var res = new(entity.User)
	res.SetUsername(in.Username)
	res.SetIsAdmin(&in.IsAdmin)
	res.SetNickname(in.Nickname)
	res.SetPassword(in.Password)
	res.SetUUID(uuid.NewV4().String())
	return res
}
