package ginmetric

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	monitor "github.com/zzzzer91/prometheus-monitor"
)

// Metric register gin metrics.
func Metric(monitor *monitor.Monitor, opts ...Option) gin.HandlerFunc {
	c := newConfig()
	c.apply(opts...)

	m := &metric{
		Monitor: monitor,
		config:  c,
	}

	m.initGinMetrics()
	return func(ctx *gin.Context) {
		startTime := time.Now()

		// execute normal process.
		ctx.Next()

		// after request
		m.ginMetricHandle(ctx, startTime)
	}
}

const (
	metricURIRequestTotal = "gin_uri_request_total"
	metricRequestBody     = "gin_request_body_total"
	metricResponseBody    = "gin_response_body_total"
	metricRequestDuration = "gin_request_duration"
	metricSlowRequest     = "gin_slow_request_total"
)

type metric struct {
	*monitor.Monitor
	*config
}

// initGinMetrics used to init gin metrics
func (m *metric) initGinMetrics() {
	m.AddMetric(&monitor.Metric{
		Type:        monitor.Counter,
		Name:        metricURIRequestTotal,
		Description: "all the server received request num with every uri.",
		Labels:      []string{"uri", "method", "code", "ip"},
	})
	m.AddMetric(&monitor.Metric{
		Type:        monitor.Counter,
		Name:        metricRequestBody,
		Description: "the server received request body size, unit byte",
		Labels:      []string{"uri", "method"},
	})
	m.AddMetric(&monitor.Metric{
		Type:        monitor.Counter,
		Name:        metricResponseBody,
		Description: "the server send response body size, unit byte",
		Labels:      []string{"uri", "method"},
	})
	m.AddMetric(&monitor.Metric{
		Type:        monitor.Histogram,
		Name:        metricRequestDuration,
		Description: "the time server took to handle the request.",
		Labels:      []string{"uri", "method"},
		Buckets:     m.reqDuration,
	})
	m.AddMetric(&monitor.Metric{
		Type:        monitor.Counter,
		Name:        metricSlowRequest,
		Description: fmt.Sprintf("the server handled slow requests counter, t=%d.", m.slowTime),
		Labels:      []string{"uri", "method"},
	})
}

func (m *metric) ginMetricHandle(ctx *gin.Context, start time.Time) {
	r := ctx.Request
	w := ctx.Writer

	labelValues := []string{ctx.FullPath(), r.Method, strconv.Itoa(w.Status()), ctx.ClientIP()}

	// set uri request total
	_ = m.GetMetric(metricURIRequestTotal).Inc(labelValues)

	// set request body size
	// since r.ContentLength can be negative (in some occasions) guard the operation
	if r.ContentLength >= 0 {
		_ = m.GetMetric(metricRequestBody).Add(labelValues[:2], float64(r.ContentLength))
	}

	// set response size
	if w.Size() > 0 {
		_ = m.GetMetric(metricResponseBody).Add(labelValues[:2], float64(w.Size()))
	}

	latency := time.Since(start)

	// set request duration
	_ = m.GetMetric(metricRequestDuration).Observe(labelValues[:2], float64(latency.Milliseconds()))

	// set slow request
	if latency > m.slowTime {
		_ = m.GetMetric(metricSlowRequest).Inc(labelValues[:2])
	}
}
