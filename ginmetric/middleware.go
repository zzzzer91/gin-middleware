package ginmetric

import (
	"github.com/gin-gonic/gin"
)

// Metrics 注册监控指标
func Metrics() gin.HandlerFunc {
	m := GetMonitor()
	m.initGinMetrics()
	return m.monitorInterceptor
}
