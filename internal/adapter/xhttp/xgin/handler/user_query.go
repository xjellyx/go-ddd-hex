package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
	"github.com/olongfen/go-ddd-hex/lib/response"
)

type userQueryCtl struct {
	domain application.UserServiceInterface
}

func (u *userQueryCtl) Get(c context.Context) {
	var (
		id  string
		res interface{}
		err error
	)
	ctx := c.(*gin.Context)
	defer func() {
		if err != nil {
			response.NewGinResponse(ctx).Fail(response.CodeFail, err).Response()
		} else {
			response.NewGinResponse(ctx).Success(res).Response()
		}
	}()

	id = ctx.GetString("id")
	data, _err := u.domain.Get(ctx.Request.Context(), vo.UserUnique{ID: id})
	if _err != nil {
		err = _err
		return
	}
	res = data

	return
}
