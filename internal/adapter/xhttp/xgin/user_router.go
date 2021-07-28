package xgin

import (
	"github.com/gin-gonic/gin"
	"github.com/olongfen/go-ddd-hex/internal/application"
)

func registerUserRouter(mux *gin.RouterGroup, userInterface application.UserInterface) {
	group := mux.Group("users")
	ctl := NewUserCtl(userInterface)
	group.GET(":id", ctl.Get)
	group.POST("/", ctl.ChangePassword)
}
