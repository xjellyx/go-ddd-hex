package xgin

import (
	"github.com/gin-contrib/pprof"
	"github.com/olongfen/go-ddd-hex/internal/application"
)

func (g *XGin) RegisterPprof() {
	pprof.Register(g.mux) // default is "debug/pprof"
}

func (g *XGin) RegisterPostRouter(post application.PostInterface) {
	group := g.mux.Group(ApiV1 + "posts")
	ctl := NewPostCtl(post)
	group.GET("/", ctl.GetByUserID)
}

func (g *XGin) RegisterUserRouter(userInterface application.UserInterface) {
	group := g.mux.Group(ApiV1 + "users")
	ctl := NewUserCtl(userInterface)
	group.GET(":id", ctl.Get)
	group.POST("/", ctl.ChangePassword)
}
