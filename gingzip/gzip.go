package gingzip

import (
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	bestSpeedGzPool = &sync.Pool{
		New: func() interface{} {
			gz, err := gzip.NewWriterLevel(io.Discard, gzip.BestSpeed)
			if err != nil {
				panic(err)
			}
			return gz
		},
	}
	bestCompressionGzPool = &sync.Pool{
		New: func() interface{} {
			gz, err := gzip.NewWriterLevel(io.Discard, gzip.BestCompression)
			if err != nil {
				panic(err)
			}
			return gz
		},
	}
	defaultCompressionGzPool = &sync.Pool{
		New: func() interface{} {
			gz, err := gzip.NewWriterLevel(io.Discard, gzip.DefaultCompression)
			if err != nil {
				panic(err)
			}
			return gz
		},
	}
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write([]byte(s))
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write(data)
}

// Fix: https://github.com/mholt/caddy/issues/38
func (g *gzipWriter) WriteHeader(code int) {
	g.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(code)
}

func Wrap(f gin.HandlerFunc, level int) gin.HandlerFunc {
	var gzPool *sync.Pool
	switch level {
	case gzip.BestSpeed:
		gzPool = bestSpeedGzPool
	case gzip.BestCompression:
		gzPool = bestCompressionGzPool
	default:
		gzPool = defaultCompressionGzPool
	}
	return func(c *gin.Context) {
		req := c.Request
		if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") &&
			!strings.Contains(req.Header.Get("Connection"), "Upgrade") &&
			!strings.Contains(req.Header.Get("Accept"), "text/event-stream") {
			gz := gzPool.Get().(*gzip.Writer)
			defer gzPool.Put(gz)
			defer gz.Reset(io.Discard)
			gz.Reset(c.Writer)

			c.Header("Content-Encoding", "gzip")
			c.Header("Vary", "Accept-Encoding")
			c.Writer = &gzipWriter{c.Writer, gz}
			defer func() {
				_ = gz.Close()
				c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
			}()
		}
		f(c)
	}
}
