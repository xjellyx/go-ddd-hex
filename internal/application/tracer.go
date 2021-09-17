package application

import (
	"github.com/opentracing/opentracing-go"
	"io"
)

type Tracer interface {
	opentracing.Tracer
	io.Closer
}
