package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/olongfen/go-ddd-hex/lib/response"
)

type postQueryCtl struct {
	domain application.PostServiceInterface
}

// GetByUserID .
// @Tags Post文章
// @Summary 获取个人文章信息默认十条
// @Description 通过用户id获取
// @Accept json
// @Produce json
// @Param userId path string true "用户id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response jwt验证失败
// @Router /api/v1/posts/:userId/ [get]
func (p *postQueryCtl) GetByUserID(ctx context.Context) {
	c := ctx.(*gin.Context)
	userId := c.Param("userId")
	var (
		res interface{}
		err error
	)
	defer func() {
		if err != nil {
			response.NewGinResponse(c).Fail(response.CodeFail, err).Response()
		} else {
			response.NewGinResponse(c).Success(res).Response()
		}
	}()
	data, _err := p.domain.GetByUserID(c.Request.Context(), userId)
	if _err != nil {
		err = _err
		return
	}

	res = data
}
