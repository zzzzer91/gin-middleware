// reference: https://github.com/joeig/gin-cachecontrol

package gincachecontrol

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const CacheControlHeader = "Cache-Control"

// Config defines a cache-control configuration.
//
// References:
// https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.2
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control
type Config struct {
	MustRevalidate       bool
	NoCache              bool
	NoStore              bool
	NoTransform          bool
	Public               bool
	Private              bool
	ProxyRevalidate      bool
	MaxAge               *time.Duration
	SMaxAge              *time.Duration
	Immutable            bool
	StaleWhileRevalidate *time.Duration
	StaleIfError         *time.Duration
}

func (c *Config) buildCacheControl() string {
	var values []string

	if c.MustRevalidate {
		values = append(values, "must-revalidate")
	}

	if c.NoCache {
		values = append(values, "no-cache")
	}

	if c.NoStore {
		values = append(values, "no-store")
	}

	if c.NoTransform {
		values = append(values, "no-transform")
	}

	if c.Public {
		values = append(values, "public")
	}

	if c.Private {
		values = append(values, "private")
	}

	if c.ProxyRevalidate {
		values = append(values, "proxy-revalidate")
	}

	if c.MaxAge != nil {
		values = append(values, fmt.Sprintf("max-age=%.f", c.MaxAge.Seconds()))
	}

	if c.SMaxAge != nil {
		values = append(values, fmt.Sprintf("s-maxage=%.f", c.SMaxAge.Seconds()))
	}

	if c.Immutable {
		values = append(values, "immutable")
	}

	if c.StaleWhileRevalidate != nil {
		values = append(values, fmt.Sprintf("stale-while-revalidate=%.f", c.StaleWhileRevalidate.Seconds()))
	}

	if c.StaleIfError != nil {
		values = append(values, fmt.Sprintf("stale-if-error=%.f", c.StaleIfError.Seconds()))
	}

	return strings.Join(values, ", ")
}

var (
	// NoCachePreset is a cache-control configuration preset which advices the HTTP client not to cache at all.
	NoCachePreset = Config{
		MustRevalidate: true,
		NoCache:        true,
		NoStore:        true,
	}

	// CacheAssetsForeverPreset is a cache-control configuration preset which advices the HTTP client
	// and all caches in between to cache the object forever without revalidation.
	// Technically, "forever" means 1 year, in order to comply with common CDN limits.
	CacheAssetsForeverPreset = Config{
		Public: true,
		MaxAge: func() *time.Duration {
			t := 8760 * time.Hour
			return &t
		}(),
		Immutable: true,
	}
)

func CacheControl(cfg *Config) gin.HandlerFunc {
	value := cfg.buildCacheControl()

	return func(c *gin.Context) {
		c.Next()
		writer := c.Writer
		if writer.Status() == http.StatusOK {
			writer.Header().Set(CacheControlHeader, value)
		}
	}
}
