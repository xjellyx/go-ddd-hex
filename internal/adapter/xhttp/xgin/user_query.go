package xgin

import (
	"github.com/gin-gonic/gin"
	"github.com/olongfen/go-ddd-hex/internal/application"

	"github.com/olongfen/go-ddd-hex/lib/response"
)

type UserQueryCtl struct {
	domain application.UserInterface
}

func (u *UserQueryCtl) Get(ctx *gin.Context) {
	var (
		id  string
		res interface{}
		err error
	)
	defer func() {
		if err != nil {
			response.NewGinResponse(ctx).Fail(response.CodeFail, err).Response()
		} else {
			response.NewGinResponse(ctx).Success(res).Response()
		}
	}()

	id = ctx.Param("id")
	data, _err := u.domain.Get(ctx.Request.Context(), id)
	if _err != nil {
		err = _err
		return
	}

	res = data

	return
}
