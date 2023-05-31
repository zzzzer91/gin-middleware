package ginmetric

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	monitor "github.com/zzzzer91/prometheus-monitor"
)

const (
	metricURIRequestTotal = "gin_uri_request_total"
	metricRequestBody     = "gin_request_body_total"
	metricResponseBody    = "gin_response_body_total"
	metricRequestDuration = "gin_request_duration"
	metricSlowRequest     = "gin_slow_request_total"
)

const (
	defaultSlowTime = 5 * time.Second
)

var (
	// Milliseconds
	defaultDuration = []float64{25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}
)

type config struct {
	m *monitor.Monitor

	slowTime    time.Duration
	reqDuration []float64
}

func newConfig(m *monitor.Monitor) *config {
	return &config{
		m:           m,
		slowTime:    defaultSlowTime,
		reqDuration: defaultDuration,
	}
}

func (c *config) apply(opts ...Option) {
	for _, o := range opts {
		o(c)
	}
}

// initGinMetrics used to init gin metrics
func (c *config) initGinMetrics() {
	c.m.AddMetric(&monitor.Metric{
		Type:        monitor.Counter,
		Name:        metricURIRequestTotal,
		Description: "all the server received request num with every uri.",
		Labels:      []string{"uri", "method", "code", "ip"},
	})
	c.m.AddMetric(&monitor.Metric{
		Type:        monitor.Counter,
		Name:        metricRequestBody,
		Description: "the server received request body size, unit byte",
		Labels:      []string{"uri", "method"},
	})
	c.m.AddMetric(&monitor.Metric{
		Type:        monitor.Counter,
		Name:        metricResponseBody,
		Description: "the server send response body size, unit byte",
		Labels:      []string{"uri", "method"},
	})
	c.m.AddMetric(&monitor.Metric{
		Type:        monitor.Histogram,
		Name:        metricRequestDuration,
		Description: "the time server took to handle the request.",
		Labels:      []string{"uri", "method"},
		Buckets:     c.reqDuration,
	})
	c.m.AddMetric(&monitor.Metric{
		Type:        monitor.Counter,
		Name:        metricSlowRequest,
		Description: fmt.Sprintf("the server handled slow requests counter, t=%d.", c.slowTime),
		Labels:      []string{"uri", "method"},
	})
}

func (c *config) ginMetricHandle(ctx *gin.Context, start time.Time) {
	r := ctx.Request
	w := ctx.Writer

	labelValues := []string{ctx.FullPath(), r.Method, strconv.Itoa(w.Status()), ctx.ClientIP()}

	// set uri request total
	c.m.GetMetric(metricURIRequestTotal).Inc(labelValues)

	// set request body size
	// since r.ContentLength can be negative (in some occasions) guard the operation
	if r.ContentLength >= 0 {
		c.m.GetMetric(metricRequestBody).Add(labelValues[:2], float64(r.ContentLength))
	}

	// set response size
	if w.Size() >= 0 {
		c.m.GetMetric(metricResponseBody).Add(labelValues[:2], float64(w.Size()))
	}

	latency := time.Since(start)

	// set request duration
	c.m.GetMetric(metricRequestDuration).Observe(labelValues[:2], float64(latency.Milliseconds()))

	// set slow request
	if latency > c.slowTime {
		c.m.GetMetric(metricSlowRequest).Inc(labelValues[:2])
	}
}

type Option func(c *config)

func WithSlowTIme(slowTime time.Duration) Option {
	return func(c *config) {
		c.slowTime = slowTime
	}
}

func WithReqDuration(reqDuration []float64) Option {
	return func(c *config) {
		c.reqDuration = reqDuration
	}
}
