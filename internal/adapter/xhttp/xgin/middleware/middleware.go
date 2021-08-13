package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
)

func GinLogFormatter() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			return fmt.Sprintf(`address: %s, time: %s, method: %s, path: %s, errMessage: %s, proto: %s, code: %d, 
latency: %s, body: %v %v`, params.ClientIP, params.TimeStamp.Format("2006-01-02 15:04:05"),
				params.Method, params.Path, params.ErrorMessage, params.Request.Proto, params.StatusCode, params.Latency,
				params.Request.Body, "\n")
		})
}

func Tracer() gin.HandlerFunc {
	return ginhttp.Middleware(opentracing.GlobalTracer())
}
