package db

import (
	"database/sql"
	"fmt"
	"github.com/olongfen/go-ddd-hex/config"
	"github.com/olongfen/go-ddd-hex/internal/contant"
	"github.com/opentracing/opentracing-go"
	tracerLog "github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	globalDB *Database
)

type Database struct {
	cfg config.DBConfig
	db  *gorm.DB
}

func NewDatabase(cfg config.DBConfig) *Database {
	globalDB = &Database{
		cfg: cfg,
	}
	return globalDB
}

func (d *Database) DB() interface{} {
	return d.db
}

func (d *Database) InjectEntities(en ...interface{}) error {
	return d.db.AutoMigrate(en...)
}

func (d *Database) Connect() (err error) {
	var (
		idb *sql.DB
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
		PrepareStmt:    true,
		NamingStrategy: schema.NamingStrategy{TablePrefix: "tb_"}, // 表明前缀
	}); err != nil {
		return
	}
	err = d.db.Use(&OpentracingPlugin{})
	if err != nil {
		return
	}
	if idb, err = d.db.DB(); err != nil {
		return
	}
	idb.SetMaxOpenConns(d.cfg.MaxOpenConns)
	idb.SetMaxIdleConns(d.cfg.MaxIdleConns)

	logrus.Infoln("db connected success")
	return
}

type OpentracingPlugin struct{}

// 告诉编译器这个结构体实现了gorm.Plugin接口
var _ gorm.Plugin = &OpentracingPlugin{}

func (op *OpentracingPlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前 - 并不是都用相同的方法，可以自己自定义
	db.Callback().Create().Before("gorm:before_create").Register(contant.CallBackBeforeName, before)
	db.Callback().Query().Before("gorm:query").Register(contant.CallBackBeforeName, before)
	db.Callback().Delete().Before("gorm:before_delete").Register(contant.CallBackBeforeName, before)
	db.Callback().Update().Before("gorm:setup_reflect_value").Register(contant.CallBackBeforeName, before)
	db.Callback().Row().Before("gorm:row").Register(contant.CallBackBeforeName, before)
	db.Callback().Raw().Before("gorm:raw").Register(contant.CallBackBeforeName, before)

	// 结束后 - 并不是都用相同的方法，可以自己自定义
	db.Callback().Create().After("gorm:after_create").Register(contant.CallBackAfterName, after)
	db.Callback().Query().After("gorm:after_query").Register(contant.CallBackAfterName, after)
	db.Callback().Delete().After("gorm:after_delete").Register(contant.CallBackAfterName, after)
	db.Callback().Update().After("gorm:after_update").Register(contant.CallBackAfterName, after)
	db.Callback().Row().After("gorm:row").Register(contant.CallBackAfterName, after)
	db.Callback().Raw().After("gorm:raw").Register(contant.CallBackAfterName, after)
	return
}

func (op *OpentracingPlugin) Name() string {
	return "opentracingPlugin"
}

func before(db *gorm.DB) {
	val, ok := db.Statement.Context.Value(contant.RepositoryMethodCtxTag).(string)
	if !ok {
		val = "gorm"
	}
	span, _ := opentracing.StartSpanFromContext(db.Statement.Context, val)
	// 利用db实例去传递span
	db.InstanceSet(contant.GormSpanKey, span)
}

func after(db *gorm.DB) {
	_span, exist := db.InstanceGet(contant.GormSpanKey)
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
