package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/olongfen/go-ddd-hex/config"
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
