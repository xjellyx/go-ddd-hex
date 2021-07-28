package xgin

import (
	"context"
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"github.com/olongfen/go-ddd-hex/config"
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/olongfen/go-ddd-hex/internal/contanst"
	"github.com/olongfen/go-ddd-hex/lib/utils"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
)

const (
	ApiV1 = "/api/v1/"
)

type XGin struct {
	appSrv map[string]interface{}
	ctx    context.Context
	cfg    *config.Config
	mux    *gin.Engine
}

func NewXGin(ctx context.Context, cfg *config.Config, app map[string]interface{}) *XGin {
	g := &XGin{
		appSrv: app,
		ctx:    ctx,
		cfg:    cfg,
		mux:    gin.Default(),
	}

	return g
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
	log.Info("http server shutdown")

}

func (g *XGin) Inject() application.XHttp {
	if !g.cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	if g.cfg.Debug {
		// 打印body
		// g.mux.Use(RequestLoggerMiddleware)
	}
	g.mux.Use(ginhttp.Middleware(g.cfg.Tracer))
	return g
}

func (g *XGin) Register() application.XHttp {
	g.RegisterPprof()
	g.RegisterPostRouter(g.appSrv[contanst.PostTag].(application.PostInterface))
	g.RegisterUserRouter(g.appSrv[contanst.UserTag].(application.UserInterface))
	return g
}
