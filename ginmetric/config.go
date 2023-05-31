package ginmetric

import "time"

const (
	defaultSlowTime = 5 * time.Second
)

var (
	// Milliseconds
	defaultDuration = []float64{25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}
)

type config struct {
	slowTime    time.Duration
	reqDuration []float64
}

func newConfig() *config {
	return &config{
		slowTime:    defaultSlowTime,
		reqDuration: defaultDuration,
	}
}

func (c *config) apply(opts ...Option) {
	for _, o := range opts {
		o(c)
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
