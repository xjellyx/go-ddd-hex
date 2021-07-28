package utils

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type ctxKeyWaitGroup struct{}

func GetWaitGroupInCtx(ctx context.Context) *sync.WaitGroup {
	if wg, ok := ctx.Value(ctxKeyWaitGroup{}).(*sync.WaitGroup); ok {
		return wg
	}

	return nil
}

func NewWaitGroupCtx() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.WithValue(context.Background(), ctxKeyWaitGroup{}, new(sync.WaitGroup)))
}

func PB2Time(in *timestamp.Timestamp) time.Time {
	if in == nil {
		return time.Time{}
	}

	r, err := ptypes.Timestamp(in)
	if err != nil {
		logrus.Errorf("err: %+v", err)
		return time.Time{}
	}

	return r
}
