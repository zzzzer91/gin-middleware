package ginmetric

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Metric register gin metrics.
func Metric(opts ...Option) gin.HandlerFunc {
	c := newConfig()
	c.apply(opts...)
	c.initGinMetrics()

	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		c.ginMetricHandle(ctx, startTime)
	}
}
