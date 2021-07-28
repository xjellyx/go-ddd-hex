package xgin

import (
	"github.com/gin-gonic/gin"
	"github.com/olongfen/go-ddd-hex/internal/application"
)

type UserExecCtl struct {
	domain application.UserInterface
}

func (u *UserExecCtl) ChangePassword(c *gin.Context) {

}
