package application

import (
	"context"
	"errors"
	"github.com/olongfen/go-ddd-hex/config"
	"github.com/olongfen/go-ddd-hex/internal/domain/dependency"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/internal/domain/service"
	"github.com/olongfen/go-ddd-hex/internal/infra/db"
	"github.com/olongfen/go-ddd-hex/internal/infra/tracer"
	"github.com/olongfen/go-ddd-hex/lib/utils"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	App = new(Application)
)

// 初始化基础组件
func init() {
	App.Ctx, App.Cancel = utils.NewWaitGroupCtx()
	App.SetDatabase(db.NewDatabase(config.GetConfig().DBConfig))
	// 数据库连接
	if err := App.db.Connect(); err != nil {
		log.Fatal(err)
	}
	if err := App.db.InjectEntities(&entity.User{}, &entity.Post{}); err != nil {
		log.Fatal(err)
	}
	// 初始化链路追踪
	App.SetTracer(tracer.GetHandlerTracer())
	opentracing.InitGlobalTracer(App.tracer)
	wg := utils.GetWaitGroupInCtx(App.Ctx)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-App.Ctx.Done()
		if err := App.tracer.Close(); err != nil {
			log.Errorln(err)
		}
		log.Infoln("Trace Close...")
	}()
}

// Application 应用程序入口
type Application struct {
	Ctx              context.Context
	Cancel           context.CancelFunc
	repos            []Repository  // 存储库
	services         []Service     // 服务
	http             XHttp         // http
	httpHandler      []IController // http 路由
	httpGroupHandler []IController // http 组路由
	db               Database      // 数据库基础组件
	tracer           Tracer
}

func checkRepo(repo Repository) (err error) {
	if repo == nil {
		return errors.New("repository is not init")
	}
	return
}

func (a *Application) GetUserService() UserServiceInterface {
	var (
		res UserServiceInterface
	)
	for _, v := range a.services {
		if val, ok := v.(UserServiceInterface); ok {
			res = val
			break
		}
	}
	if res == nil {
		log.Fatal("user service not append")
	}
	return res
}

// InjectServices 注册服务
func (a *Application) InjectServices() *Application {
	var (
		err         error
		userRepo    dependency.UserRepo
		postRepo    dependency.PostRepo
		userService UserServiceInterface
		postService PostServiceInterface
	)
	for _, v := range a.repos {
		switch v.(type) {
		case dependency.UserRepo:
			userRepo = v.(dependency.UserRepo)
		case dependency.PostRepo:
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
	userService = service.NewUserService(userRepo)
	postService = service.NewPostService(postRepo, userRepo)
	a.services = append(a.services, userService, postService)
	for _, v := range a.httpGroupHandler {
		switch v.(type) {
		case PostHandler:
			v.(PostHandler).SetService(postService)
		case UserHandler:
			v.(UserHandler).SetService(userService)
		}
	}
	a.http.Init().Use()
	return a
}

// AppendRepo 添加存储库
func (a *Application) AppendRepo(repo Repository) *Application {
	a.repos = append(a.repos, repo)
	return a
}

// AppendHTTPHandler http handler
func (a *Application) AppendHTTPHandler(ctl IController) *Application {
	a.httpHandler = append(a.httpHandler, ctl)
	return a
}

// AppendHTTPGroupHandler http handler
func (a *Application) AppendHTTPGroupHandler(ctl IController) *Application {
	a.httpGroupHandler = append(a.httpGroupHandler, ctl)
	return a
}

// SetHttp 设置http适配器
func (a *Application) SetHttp(x XHttp) *Application {
	a.http = x
	return a
}

// SetDatabase 设置数据库基础组件
func (a *Application) SetDatabase(d Database) *Application {
	a.db = d
	return a
}

// SetTracer 设置链路追踪基础组件
func (a *Application) SetTracer(t Tracer) *Application {
	a.tracer = t
	return a
}

// GetDB 设置数据库基础组件
func (a *Application) GetDB() Database {
	return a.db
}

func (a *Application) Run() {
	go App.http.RouterGroup("api/v1", App.httpGroupHandler...).Route(App.httpHandler...).Run()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	var state int32 = 1
EXIT:
	for {
		sig := <-quit
		App.Cancel()
		utils.GetWaitGroupInCtx(App.Ctx).Wait()
		log.Printf("signal[%s]\n", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			atomic.StoreInt32(&state, 0)
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	log.Println("Program Exit...")
	time.Sleep(time.Second)
	os.Exit(int(atomic.LoadInt32(&state)))
}
