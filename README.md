# golang 领域驱动设计（六边形）

- 六边形领域驱动
```
    六边形领域驱动模型也称为端口适配器模型，它分为外部和内部，内部包含application和domain层，外部包含适配器层，依赖关系为
adapter->aplication->domain,domain层业务代码只向application层暴露。
   domain领域层：
       基本包含实体（entity）、值对象value object（vo）、聚合（aggregate）、依赖倒置（denpendency）、服务（sevice）实现业务逻辑（domain层
    里面如需其他业务可以自行扩展）。
   application应用层：
       没有太多逻辑，定义接口适配对象。
   adapter适配器层: 
       主适配器（例如 http）和从适配器（例如database）。
   （个人对于六边形架构只是进行简单的描述，如需详细描述） 
```
[Netflix 的六边形架构实践](https://www.infoq.cn/article/pjekymkzhmkafgi6ycri)

[Domain-Driven Design and the Hexagonal Architecture](https://vaadin.com/learn/tutorials/ddd/ddd_and_hexagonal)

- 项目实例结构
```
  │ ├─config             项目配置
  │ ├─internal           项目核心代码
        │ ├─adapter      适配器
               │ ├─xhttp http适配器
               │ ├─respository 存储适配器      
        │ ├─application  应用服务层
        │ ├─contanst     常量
        │ ├─domain       领域层
              │ ├─aggregate 聚合函数
              │ ├─dependency 依赖倒置
              │ ├─entity 实体
              │ ├─service 业务逻辑
              │ ├─vo      值对象
        │ ├─infra         基础组件   
  │ ├─lib                (可以公共的工具)
```

- main入口

    ```go
            package main
            
            import (
            	_ "github.com/olongfen/go-ddd-hex/internal/adapter/repository" // 初始化存储库适配器
            	_ "github.com/olongfen/go-ddd-hex/internal/adapter/xhttp/xgin" // 初始化http适配器
            	"github.com/olongfen/go-ddd-hex/internal/application"
            )
            
            func main() {
            	application.App.InjectServices().Run()
            }
    ```

- application应用层主要代码
    ```go
        package application
        
        import (
        	"context"
        	"errors"
        	"github.com/olongfen/go-ddd-hex/config"
        	"github.com/olongfen/go-ddd-hex/internal/domain/aggregate"
        	"github.com/olongfen/go-ddd-hex/internal/domain/dependency"
        	"github.com/olongfen/go-ddd-hex/internal/domain/service"
        	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
        	"github.com/olongfen/go-ddd-hex/internal/infra/db"
        	"github.com/olongfen/go-ddd-hex/lib/utils"
        	"github.com/opentracing/opentracing-go"
        	log "github.com/sirupsen/logrus"
        	"github.com/uber/jaeger-client-go"
        	jaegercfg "github.com/uber/jaeger-client-go/config"
        	"github.com/uber/jaeger-client-go/log/zap"
        	"github.com/uber/jaeger-lib/metrics"
        	"io"
        	"os"
        	"reflect"
        )
        
        var (
        	App = new(Application)
        )
        
        func init() {
        	App.Ctx, App.Cancel = utils.NewWaitGroupCtx()
        	cfg := config.GetConfig()
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
    ```

- domain领域层

    - entity 实体
        - user.go用户实体
        ```go
        type User struct {
            gorm.Model
            UUID     string      `gorm:"uniqueIndex;not null;type:varchar(36)"`
            Username string      `gorm:"uniqueIndex;not null;type:varchar(36)"` // 用户名
            Password null.String `gorm:"type:varchar(16)"`                      // 密码
            Nickname null.String `gorm:"type:varchar(36)"`                      // 昵称
            IsAdmin  null.Bool   `gorm:"default: false"`                        // true：是管理员
        }
      ```

        - post.go帖子实体
        ```go
            type Post struct {
                gorm.Model
                UserUUID string `gorm:"not null;type:varchar(36)"`
                Title    string `gorm:"type:varchar(64)"` // 文章标题
                Content  string // 文章内容
            }
      ```  
    - vo  值对象
        - user.go 用户值对象
            ```go
                type UserRes struct {
                    BaseRes
                    UUID     string `json:"uuid,omitempty"`
                    Username string `json:"username,omitempty"` // 用户名
                    Nickname string `json:"nickname,omitempty"` // 昵称
                    IsAdmin  bool   `json:"isAdmin,omitempty"`  // true：是管理员
                     }
          ``` 
        - post 帖子值对象
            ```go
              type PostRes struct {
                BaseRes
                Title   string `json:"title"`
                Content string `json:"content"`
                }  
          ``` 
        - base 公共
            ```go
                type BaseRes struct {
                    ID        string    `json:"id,omitempty"`
                    CreatedAt time.Time `json:"createdAt,omitempty"`
                    UpdatedAt time.Time `json:"updatedAt,omitempty"`
                    }
          ```
    -   aggregate 聚合
        ```go

          // UserPostFactory 帖子与用户逻辑聚合工厂模式
            type UserPostFactory struct {
                UserRepo dependency.UserRepo // 用户存储库
                PostRepo dependency.PostRepo // 帖子存储库
            }
            
            type QueryUserPostRes struct {
                User  vo.UserRes   `json:"user"`
                Posts []vo.PostRes `json:"posts"`
            }
        ```
    - dependency 依赖倒置
        - user.go 用户存储对象
           ```go
           type UserRepo interface {
               Get(id string) (*entity.User, error)
               Find(cond map[string]interface{}, meta *query.Meta) ([]*entity.User, error)
               Create(user *entity.User) error
               Update(cond map[string]interface{}, change interface{}) error
               Delete(cond map[string]interface{}) error
                }
          ```
        - post.go 帖子存储对象
          ```go
               type PostRepo interface {
                   Get(id string) (*entity.Post, error)
                   Find(cond map[string]interface{}, meta *query.Meta) ([]*entity.Post, error)
                   Create(post *entity.Post) error
                   Update(cond map[string]interface{}, change interface{}) error
                   Delete(cond map[string]interface{}) error
               }
          ```
    - service 业务逻辑

- adapter 适配器层
    - repository 实现dependency定义的接口
    - xhttp 实现application应用层定义的xhttp接口

# 总结
        domain层的代码不能依赖外面任何一层，只让application层可以依赖使用，adapter层依赖应用层而且实现出入端口具体逻辑。
    面向对象开发是使用领域驱动框架必不可少的一部分，六边形模型可以解决很多让人蛋疼的业务开发，如果业务逻辑不是很复杂的话使用
    六边形框架设计起来会比较麻烦。 另外ddd清洁分层架构也是不错的框架， 领域驱动的原理和详细google有很多大佬详写，大家需要看
    详细内容可以自行google或者阅读上面我推荐两个链接文章，所以一些具体的解释本人就不献丑了，本项目是本人学习领域驱动框架之后
    进行设想开发的，欢迎大家一起指正或者改善。个人觉得应用层可以设计层框架使用，目前暂时还没有idea，也欢迎大佬们讨论。
[项目地址](https://github.com/olongfen/go-ddd-hex)