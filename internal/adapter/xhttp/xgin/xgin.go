package xgin

import (
	"context"
	"crypto/tls"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/olongfen/go-ddd-hex/config"
	_ "github.com/olongfen/go-ddd-hex/docs"
	"github.com/olongfen/go-ddd-hex/internal/adapter/xhttp/xgin/middleware"
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/olongfen/go-ddd-hex/lib/utils"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net"
	"net/http"
	"time"
)

var (
	_ application.XHttp = (*XGin)(nil)
)

const (
	ApiV1 = "/api/v1/"
)

type XGin struct {
	ctx context.Context
	cfg *config.Config
	mux *gin.Engine
}

func init() {
	g := &XGin{
		ctx: application.App.Ctx,
		cfg: config.GetConfig(),
		mux: gin.Default(),
	}
	application.App.SetHttp(g)
}

func (g *XGin) Run() {
	server := &http.Server{
		Addr:    net.JoinHostPort("", g.cfg.HttpPort),
		Handler: g.mux,
	}

	var startFn func()
	if g.cfg.TlsCert != "" && g.cfg.TlsKey != "" {
		// https
		startFn = func() {
			log.Infof("https server start: %v", server.Addr)
			cer, err := tls.LoadX509KeyPair(g.cfg.TlsCert, g.cfg.TlsKey)
			if err != nil {
				log.Errorf("failed to load certificate and key: %v", err)
				return
			}
			tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
			server.TLSConfig = tlsConfig

			if err := server.ListenAndServeTLS(g.cfg.TlsCert, g.cfg.TlsKey); err != nil && err != http.ErrServerClosed {
				log.Fatalf("err: %+v", err)
			}
		}
	} else {
		// http
		startFn = func() {
			log.Infof("http server start: %v", server.Addr)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("err: %+v", err)
			}
		}
	}

	go startFn()
	wg := utils.GetWaitGroupInCtx(g.ctx)
	wg.Add(1)
	defer wg.Done()
	<-g.ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("server shutdown err: %+v", err)
		return
	}
	log.Info("HTTP Server Shutdown...")
}

func (g *XGin) Register(repos []application.Service) application.XHttp {
	g.RegisterPprof()
	if !g.cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	if g.cfg.Debug {
		// 打印body
		g.mux.Use(middleware.GinLogFormatter())
	}
	// 使用中间件
	g.mux.Use(cors.Default())
	g.mux.GET("swagger/*any",ginSwagger.WrapHandler(swaggerFiles.Handler))
	g.mux.Use(middleware.Tracer())
	for _, v := range repos {
		switch v.(type) {
		case application.UserInterface:
			g.registerUserRouter(v.(application.UserInterface))
		case application.PostInterface:
			g.registerPostRouter(v.(application.PostInterface))
		}
	}

	return g
}
