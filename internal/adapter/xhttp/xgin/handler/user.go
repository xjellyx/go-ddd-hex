package handler

import (
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/sirupsen/logrus"
)

var (
	_ application.UserHandler = (*userCtl)(nil)
)

func init() {
	application.App.AppendHTTPGroupHandler(&userCtl{})
}

type userCtl struct {
	userExecCtl
	userQueryCtl
}

func (u *userCtl) SetService(domain application.UserServiceInterface) {
	u.userQueryCtl.domain = domain
	u.userExecCtl.domain = domain
}

func (u *userCtl) Router(xHttp application.XHttp, isGroup bool) {
	if u.userQueryCtl.domain == nil || u.userExecCtl.domain == nil {
		logrus.Fatal("user controller domain service not init")
	}
	const user = "users"
	xHttp.Handle("GET", user+"/:id", u.Get, isGroup)
	xHttp.Handle("PUT", user+"/change-passwd", u.ChangePasswd, isGroup)
	xHttp.Handle("POST", user+"/", u.Create, isGroup)
}
