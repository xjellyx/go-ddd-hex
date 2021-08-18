package vo

import (
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	uuid "github.com/satori/go.uuid"
)

type UserVO struct {
	BaseVO
	UUID     string `json:"uuid,omitempty"`
	Username string `json:"username,omitempty"` // 用户名
	Nickname string `json:"nickname,omitempty"` // 昵称
	IsAdmin  bool   `json:"isAdmin,omitempty"`  // true：是管理员
}

type UserVOForm struct {
	Username string `json:"username,omitempty"` // 用户名
	Nickname string `json:"nickname,omitempty"` // 昵称
	IsAdmin  bool   `json:"isAdmin,omitempty"`  // true：是管理员
	Password string `json:"password"`
}

func UserEntity2VO(in *entity.User) *UserVO {
	res := new(UserVO)
	res.Username = in.Username
	res.Nickname = in.Nickname.String
	res.IsAdmin = in.IsAdmin.Bool
	res.CreatedAt = in.CreatedAt
	res.UpdatedAt = in.UpdatedAt
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
