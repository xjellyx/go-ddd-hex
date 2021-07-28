package xgin

import (
	"github.com/gin-gonic/gin"
	"github.com/olongfen/go-ddd-hex/internal/application"
)

func registerPostRouter(mux *gin.RouterGroup, post application.PostInterface) {
	group := mux.Group("posts")
	ctl := NewPostCtl(post)
	group.GET("/", ctl.GetByUserID)
}
