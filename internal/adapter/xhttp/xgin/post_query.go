package xgin

import (
	"github.com/gin-gonic/gin"
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/olongfen/go-ddd-hex/internal/domain/aggregate"
	"github.com/olongfen/go-ddd-hex/lib/response"
)

type PostQueryCtl struct {
	domain application.PostInterface
}

func (p *PostQueryCtl) GetByUserID(c *gin.Context) {
	userId := c.Query("userId")
	var (
		data *aggregate.QueryUserPostRes
		err  error
	)
	defer func() {
		if err != nil {
			response.NewGinResponse(c).Fail(response.CodeFail, err).Response()
		} else {
			response.NewGinResponse(c).Success(data).Response()
		}
	}()
	if data, err = p.domain.GetByUserID(userId); err != nil {
		return
	}
}
