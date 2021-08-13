package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/olongfen/go-ddd-hex/config"
	"github.com/opentracing/opentracing-go"
	tracerLog "github.com/opentracing/opentracing-go/log"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"reflect"
)

var (
	globalDB  *Database
	injectors []func(db *gorm.DB)
)

type Database struct {
	cfg *config.DBConfig
	db  *gorm.DB
}

func NewDatabase(cfg *config.DBConfig) *Database {
	globalDB = &Database{
		cfg: cfg,
	}
	return globalDB
}

func (d *Database) DB(ctx context.Context) interface{} {
	iface := ctx.Value(ctxTransactionKey{})

	if iface != nil {
		tx, ok := iface.(*gorm.DB)
		if !ok {
			logrus.Panicf("unexpect context value type: %s", reflect.TypeOf(tx))
			return nil
		}

		return tx
	}

	return d.db.WithContext(ctx)
}

func (d *Database) Connect() {
	var (
		idb *sql.DB
		err error
	)
	dsn := fmt.Sprintf(`%s://%s:%s@%s:%s/%s?sslmode=disable`, d.cfg.Driver, d.cfg.Username, d.cfg.Password, d.cfg.Host,
		d.cfg.Port, d.cfg.DatabaseName)
	var ormLogger logger.Interface
	if d.cfg.Debug {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}

	if d.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         ormLogger,
		NamingStrategy: schema.NamingStrategy{TablePrefix: "tb_"}, // 表明前缀
	}); err != nil {
		logrus.Fatal(err)
	}
	d.db.Use(&OpentracingPlugin{})
	if idb, err = d.db.DB(); err != nil {
		logrus.Fatal(err)
	}
	idb.SetMaxOpenConns(d.cfg.MaxOpenConns)
	idb.SetMaxIdleConns(d.cfg.MaxIdleConns)

	registerCallback(d.db)
	callInjector(d.db)

	logrus.Infoln("db connected success")
}

// RegisterInjector 注册注入器
func RegisterInjector(f func(*gorm.DB)) {
	injectors = append(injectors, f)
}

func callInjector(db *gorm.DB) {
	for _, v := range injectors {
		v(db)
	}
}

func registerCallback(db *gorm.DB) {
	// 自动添加uuid
	err := db.Callback().Create().Before("gorm:create").Register("uuid", func(db *gorm.DB) {
		db.Statement.SetColumn("id", uuid.NewV4().String())
	})
	if err != nil {
		logrus.Panicf("err: %+v", err)
	}
}

// 包内静态变量
const gormSpanKey = "__gorm_span"
const (
	RepositoryMethodCtxTag = "repository_method"
	callBackBeforeName     = "opentracing:before"
	callBackAfterName      = "opentracing:after"
)

type OpentracingPlugin struct{}

// 告诉编译器这个结构体实现了gorm.Plugin接口
var _ gorm.Plugin = &OpentracingPlugin{}

func (op *OpentracingPlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前 - 并不是都用相同的方法，可以自己自定义
	db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)

	// 结束后 - 并不是都用相同的方法，可以自己自定义
	db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)
	return
}

func (op *OpentracingPlugin) Name() string {
	return "opentracingPlugin"
}

func before(db *gorm.DB) {
	val, ok := db.Statement.Context.Value(RepositoryMethodCtxTag).(string)
	if !ok {
		val = "gorm"
	}
	span, _ := opentracing.StartSpanFromContext(db.Statement.Context, val)
	// 利用db实例去传递span
	db.InstanceSet(gormSpanKey, span)
}

func after(db *gorm.DB) {
	_span, exist := db.InstanceGet(gormSpanKey)
	if !exist {
		return
	}
	// 断言类型
	span, ok := _span.(opentracing.Span)
	if !ok {
		return
	}

	defer span.Finish()

	if db.Error != nil {
		span.LogFields(tracerLog.Error(db.Error))
	}

	span.LogFields(tracerLog.String("sql", db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)))

}
