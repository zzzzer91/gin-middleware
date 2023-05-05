package ginmetric

import (
	"github.com/gin-gonic/gin"
)

// Metrics register monitoring metrics.
func Metrics() gin.HandlerFunc {
	m := GetMonitor()
	m.initGinMetrics()
	return m.monitorInterceptor
}
