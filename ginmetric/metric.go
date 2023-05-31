package ginmetric

import (
	"time"

	"github.com/gin-gonic/gin"
	monitor "github.com/zzzzer91/prometheus-monitor"
)

// Metric register gin metrics.
func Metric(m *monitor.Monitor, opts ...Option) gin.HandlerFunc {
	c := newConfig(m)
	c.apply(opts...)
	c.initGinMetrics()

	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		c.ginMetricHandle(ctx, startTime)
	}
}
