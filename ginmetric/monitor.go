package ginmetric

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

// Monitor is an object that uses to set gin server monitor.
type Monitor struct {
	slowTime    int32
	reqDuration []float64
	metrics     map[string]*Metric
}

const (
	metricURIRequestTotal = "gin_uri_request_total"
	metricRequestBody     = "gin_request_body_total"
	metricResponseBody    = "gin_response_body_total"
	metricRequestDuration = "gin_request_duration"
	metricSlowRequest     = "gin_slow_request_total"
)

var (
	defaultDuration = []float64{0.05, 0.1, 0.3, 1.2, 5, 10}
	monitor         *Monitor
)

// GetMonitor used to get global Monitor object,
// this function returns a singleton object.
func GetMonitor() *Monitor {
	if monitor == nil {
		monitor = &Monitor{
			slowTime:    defaultSlowTime,
			reqDuration: defaultDuration,
			metrics:     make(map[string]*Metric),
		}
	}
	return monitor
}

// GetMetric used to get metric object by metric_name.
func (m *Monitor) GetMetric(name string) *Metric {
	if metric, ok := m.metrics[name]; ok {
		return metric
	}
	return &Metric{}
}

// SetSlowTime set slowTime property. slowTime is used to determine whether
// the request is slow. For "gin_slow_request_total" metric.
func (m *Monitor) SetSlowTime(slowTime int32) {
	m.slowTime = slowTime
}

// SetDuration set reqDuration property. reqDuration is used to ginRequestDuration
// metric buckets.
func (m *Monitor) SetDuration(duration []float64) {
	m.reqDuration = duration
}

// AddMetric add custom monitor metric.
func (m *Monitor) AddMetric(metric *Metric) error {
	if _, ok := m.metrics[metric.Name]; ok {
		return errors.Errorf("metric '%s' is existed", metric.Name)
	}

	if metric.Name == "" {
		return errors.Errorf("metric name cannot be empty.")
	}
	if f, ok := promTypeHandler[metric.Type]; ok {
		if err := f(metric); err == nil {
			prometheus.MustRegister(metric.vec)
			m.metrics[metric.Name] = metric
			return nil
		}
	}
	return errors.Errorf("metric type '%d' not existed.", metric.Type)
}

// initGinMetrics used to init gin metrics
func (m *Monitor) initGinMetrics() {
	_ = monitor.AddMetric(&Metric{
		Type:        Counter,
		Name:        metricURIRequestTotal,
		Description: "all the server received request num with every uri.",
		Labels:      []string{"uri", "method", "code"},
	})
	_ = monitor.AddMetric(&Metric{
		Type:        Counter,
		Name:        metricRequestBody,
		Description: "the server received request body size, unit byte",
		Labels:      []string{"uri", "method"},
	})
	_ = monitor.AddMetric(&Metric{
		Type:        Counter,
		Name:        metricResponseBody,
		Description: "the server send response body size, unit byte",
		Labels:      []string{"uri", "method"},
	})
	_ = monitor.AddMetric(&Metric{
		Type:        Histogram,
		Name:        metricRequestDuration,
		Description: "the time server took to handle the request.",
		Labels:      []string{"uri", "method"},
		Buckets:     m.reqDuration,
	})
	_ = monitor.AddMetric(&Metric{
		Type:        Counter,
		Name:        metricSlowRequest,
		Description: fmt.Sprintf("the server handled slow requests counter, t=%d.", m.slowTime),
		Labels:      []string{"uri", "method"},
	})
}

// monitorInterceptor as gin monitor middleware.
func (m *Monitor) monitorInterceptor(ctx *gin.Context) {
	startTime := time.Now()

	// execute normal process.
	ctx.Next()

	// after request
	m.ginMetricHandle(ctx, startTime)
}

func (m *Monitor) ginMetricHandle(ctx *gin.Context, start time.Time) {
	r := ctx.Request
	w := ctx.Writer

	labelValues := []string{ctx.FullPath(), r.Method, strconv.Itoa(w.Status())}

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
	_ = m.GetMetric(metricRequestDuration).Observe(labelValues[:2], latency.Seconds())

	// set slow request
	if int32(latency.Seconds()) > m.slowTime {
		_ = m.GetMetric(metricSlowRequest).Inc(labelValues[:2])
	}
}
