package config

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/olongfen/go-ddd-hex/lib/utils"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/log/zap"
	"github.com/uber/jaeger-lib/metrics"
	"io"
	"io/fs"
	"os"
)

const (
	runModeUnitTest = "unit-test"
)

var (
	cfg = new(Config)
)

type Config struct {
	APPName         string
	HttpPort        string
	GrpcGatewayPort string
	GrpcPort        string
	Debug           bool // debug log
	// https
	TlsCert string
	TlsKey  string

	// db
	DBConfig
	Ctx    context.Context
	Cancel context.CancelFunc
	/*
		根据环境变量配置jaeger，参考 https://github.com/jaegertracing/jaeger-client-go#environment-variables

		JAEGER_AGENT_HOST
		JAEGER_AGENT_PORT
	*/
	Tracer opentracing.Tracer
}

type DBConfig struct {
	Host         string
	Port         string
	Driver       string
	DatabaseName string
	Username     string
	Password     string
	MaxIdleConns int
	MaxOpenConns int
	AutoMigrate  bool // 自动建表，补全缺失字段，初始化数据
	Debug        bool
}

func setDefault() {
	cfg.Ctx, cfg.Cancel = utils.NewWaitGroupCtx()
	viper.SetDefault("appname", "user_server")
	viper.SetDefault("httpPort", "8100")
	viper.SetDefault("grpcPort", "8200")
	viper.SetDefault("grpcGatewayPort", "8300")
	viper.SetDefault("debug", true)

	// set database default
	viper.SetDefault("dbconfig.host", "127.0.0.1")
	viper.SetDefault("dbconfig.port", "5432")
	viper.SetDefault("dbconfig.Driver", "postgres")
	viper.SetDefault("dbconfig.DatabaseName", "postgres")
	viper.SetDefault("dbconfig.Username", "postgres")
	viper.SetDefault("dbconfig.Password", "123456")
	viper.SetDefault("dbconfig.MaxIdleConns", "10")
	viper.SetDefault("dbconfig.MaxOpenConns", "20")
	viper.SetDefault("dbconfig.AutoMigrate", true)
	viper.SetDefault("dbconfig.Debug", true)
}

func setTrace(ctx context.Context) (err error) {
	defer func() {
		cfg.Tracer = opentracing.GlobalTracer()
	}()
	isSet := func(env string) bool {
		_, ok := os.LookupEnv(env)
		return ok
	}

	if !(isSet("JAEGER_AGENT_HOST") ||
		isSet("JAEGER_ENDPOINT")) {
		return
	}
	var (
		jaegerCfg *jaegercfg.Configuration
		closer    io.Closer
	)
	if jaegerCfg, err = jaegercfg.FromEnv(); err != nil {
		return
	}
	if cfg.Debug {
		jaegerCfg.Sampler.Type = jaeger.SamplerTypeConst
		jaegerCfg.Sampler.Param = 1
		jaegerCfg.Reporter.LogSpans = true
	}
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	if closer, err = jaegerCfg.InitGlobalTracer(cfg.APPName, jaegercfg.Logger(zap.NewLogger(nil)),
		jaegercfg.Metrics(jMetricsFactory)); err != nil {
		log.Errorf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	wg := utils.GetWaitGroupInCtx(ctx)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err = closer.Close(); err != nil {
			log.Errorln(err)
		}
		log.Infoln("trace close")
	}()
	return
}

func init() {
	setDefault()
	viper.SetEnvPrefix("user")
	// 运行模式,单元测试需要输入配置文件目录
	_ = viper.BindEnv("run_mode")
	_ = viper.BindEnv("config_dir")
	runMode := viper.Get("run_mode")

	configDir := viper.Get("config_dir")
	if runMode == runModeUnitTest && configDir == nil {
		log.Fatal("单元测试模式请输入配置文件绝对路径")
	}
	if configDir == nil {
		configDir = "."
	}
	viper.SetConfigFile(fmt.Sprintf(`%s%s`, configDir, "/config/config.yaml"))
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*fs.PathError); ok {
			// 配置文件未找到错误；如果需要可以忽略
		} else {
			log.Fatal(err)
		}
	}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatal(err)
	}

	if err := setTrace(cfg.Ctx); err != nil {
		log.Fatal(err)
	}

	if err := viper.WriteConfig(); err != nil {
		log.Fatal(err)
	}
	// 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Infoln("配置文件修改更新!")
		if err := viper.Unmarshal(cfg); err != nil {
			log.Fatal(err)
		}
	})

}

func GetConfig() *Config {
	return cfg
}
