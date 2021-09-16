package middleware

import (
	"context"
	"dcs-sched/pkg/common"
	"github.com/carter115/gslog"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				gslog.Errorf(context.Background(), "panic: %+v", err)
				c.JSON(200, common.NewResponse(c, common.Unknown, nil))
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}
