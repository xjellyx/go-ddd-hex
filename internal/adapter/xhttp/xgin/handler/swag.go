package handler

import (
	"github.com/gin-gonic/gin"
	_ "github.com/olongfen/go-ddd-hex/docs"
	"github.com/olongfen/go-ddd-hex/internal/application"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type swagCtl struct {
}

func init() {
	application.App.AppendHTTPHandler(&swagCtl{})
}

func (s *swagCtl) Router(xhttp application.XHttp, isGroup bool) {
	xhttp.GetMux().(*gin.Engine).GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
