package utils

import (
	"context"
	"github.com/gin-gonic/gin"
)

func WrapF(f func(ctx context.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(c)
	}
}
