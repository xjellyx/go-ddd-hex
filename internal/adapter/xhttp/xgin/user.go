package xgin

import (
	"github.com/olongfen/go-ddd-hex/internal/application"
)

type UserCtl struct {
	UserExecCtl
	UserQueryCtl
}

func NewUserCtl(domain application.UserInterface) *UserCtl {
	return &UserCtl{UserExecCtl{domain: domain}, UserQueryCtl{domain: domain}}
}
