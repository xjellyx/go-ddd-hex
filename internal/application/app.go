package application

import (
	"context"
	"errors"
	"github.com/olongfen/go-ddd-hex/config"
	"github.com/olongfen/go-ddd-hex/internal/domain/aggregate"
	"github.com/olongfen/go-ddd-hex/internal/domain/dependency"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/internal/domain/service"
	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
	"github.com/olongfen/go-ddd-hex/internal/infra/db"
	"github.com/olongfen/go-ddd-hex/lib/utils"
	"github.com/opentracing/opentracing-go"
	prometheus "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"
	"gorm.io/gorm"
	"io"
	"reflect"
)

var (
	App = new(Application)
)

func init() {
	App.Ctx, App.Cancel = utils.NewWaitGroupCtx()
	cfg := config.GetConfig()
	db.RegisterInjector(func(db *gorm.DB) {
		if config.GetConfig().AutoMigrate {
			err := db.AutoMigrate(&entity.User{}, &entity.Post{})
			if err != nil {
				log.Fatal(err)
			}
		}
	})
	// 数据库初始化
	App.SetDatabase(db.NewDatabase(&cfg.DBConfig))
	App.Connect()
	if err := App.setTrace(); err != nil {
		log.Fatal(err)
	}
}

// UserInterface user 用户服务接口
type UserInterface interface {
	ChangePassword(id string, oldPwd, newPwd string) error
	Get(id string) (res *vo.UserRes, err error)
}

// PostInterface post 服务接口
type PostInterface interface {
	GetByUserID(userID string) (*aggregate.QueryUserPostRes, error)
}

// XHttp  http 接口
type XHttp interface {
	Run()
	Register(reps []Service) XHttp
}

// Database 数据库基础组件接口
type Database interface {
	Connect()
	DB(ctx context.Context) interface{}
}

// Repository 存储库接口
type Repository interface{}

// Service service 服务接口
type Service interface {
}

// Application 应用程序入口
type Application struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	/*
		根据环境变量配置jaeger，参考 https://github.com/jaegertracing/jaeger-client-go#environment-variables

		JAEGER_AGENT_HOST
		JAEGER_AGENT_PORT
	*/
	Tracer   opentracing.Tracer
	repos    []Repository
	services []Service
	XHttp
	Database
}

func checkRepo(repo Repository) (err error) {
	if repo == nil {
		return errors.New("repository is not init")
	}
	return
}

// InjectServices 注册服务
func (a *Application) InjectServices() *Application {
	var (
		err         error
		userRepo    dependency.UserRepo
		postRepo    dependency.PostRepo
		userService UserInterface
		postService PostInterface
	)
	for _, v := range a.repos {
		t := reflect.TypeOf(v)
		switch t.Elem().Name() {
		case "userRepo":
			userRepo = v.(dependency.UserRepo)
		case "postRepo":

			postRepo = v.(dependency.PostRepo)

		}
	}
	// 验证存储库是否已经初始化
	if err = checkRepo(userRepo); err != nil {
		log.Fatal(err)
	}
	if err = checkRepo(postRepo); err != nil {
		log.Fatal(err)
	}
	// 注册服务
	userService = service.NewUserService(db.NewTxImpl(), userRepo)
	postService = service.NewPostService(db.NewTxImpl(), postRepo, userRepo)
	a.services = append(a.services, userService)
	a.services = append(a.services, postService)
	a.Register(a.services)
	return a
}

// AppendRepo 添加存储库
func (a *Application) AppendRepo(repo Repository) *Application {
	a.repos = append(a.repos, repo)
	return a
}

// SetXHttp 设置http适配器
func (a *Application) SetXHttp(x XHttp) *Application {
	a.XHttp = x
	return a
}

// SetDatabase 设置数据库基础组件
func (a *Application) SetDatabase(d Database) *Application {
	a.Database = d
	return a
}

func (a *Application) setTrace() (err error) {
	cfg := config.GetConfig()
	var (
		closer io.Closer
	)
	jaegerCfg := &jaegercfg.Configuration{
		ServiceName: cfg.APPName,
		Reporter: &jaegercfg.ReporterConfig{LogSpans: true,
			CollectorEndpoint: cfg.JaegerEndpoint,
		},
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1}}
	jMetricsFactory := jprom.New(jprom.WithRegisterer(prometheus.NewPedanticRegistry()))
	jLogger := jaegerlog.StdLogger
	if a.Tracer, closer, err = jaegerCfg.NewTracer(jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory)); err != nil {
		log.Fatal(err)
	}
	wg := utils.GetWaitGroupInCtx(a.Ctx)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-a.Ctx.Done()
		if err = closer.Close(); err != nil {
			log.Errorln(err)
		}
		log.Infoln("trace close")
	}()
	return
}
