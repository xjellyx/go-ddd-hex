package vo

import "github.com/olongfen/go-ddd-hex/internal/domain/entity"

type UserVO struct {
	BaseVO
	UUID     string `json:"uuid,omitempty"`
	Username string `json:"username,omitempty"` // 用户名
	Nickname string `json:"nickname,omitempty"` // 昵称
	IsAdmin  bool   `json:"isAdmin,omitempty"`  // true：是管理员
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
