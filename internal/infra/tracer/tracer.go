package tracer

import (
	"github.com/olongfen/go-ddd-hex/config"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"
	"io"
)

type HandlerTracer struct {
	/*
		根据环境变量配置jaeger，参考 https://github.com/jaegertracing/jaeger-client-go#environment-variables

		JAEGER_AGENT_HOST
		JAEGER_AGENT_PORT
	*/
	opentracing.Tracer
	io.Closer
}

var (
	globalHandlerTracer *HandlerTracer
)

func GetHandlerTracer() *HandlerTracer {
	return globalHandlerTracer
}

func init() {
	cfg := config.GetConfig()
	var (
		err    error
		tracer opentracing.Tracer
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
	httpHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "http_server",
		Subsystem:   "",
		Name:        "requests_seconds",
		Help:        "Histogram of response latency (seconds) of http handlers.",
		ConstLabels: nil,
		Buckets:     nil,
	}, []string{"method", "code", "uri"})
	prometheus.MustRegister(httpHistogram)
	jMetricsFactory := jprom.New(jprom.WithRegisterer(prometheus.DefaultRegisterer))
	jLogger := jaegerlog.StdLogger
	if tracer, closer, err = jaegerCfg.NewTracer(jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory)); err != nil {
		log.Fatal(err)
	}

	globalHandlerTracer = &HandlerTracer{
		Tracer: tracer, Closer: closer,
	}

}
