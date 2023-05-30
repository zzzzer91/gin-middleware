package gintrace

import (
	"github.com/gin-gonic/gin"
	"github.com/zzzzer91/zlog/trace"
	"go.opentelemetry.io/otel/attribute"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := trace.StartTracing(c.Request.Context(), c.Request.Method+" "+c.FullPath())
		if span.IsRecording() {
			defer span.End()
			path := c.Request.URL.Path
			if len(c.Request.URL.RawQuery) > 0 {
				path += "?" + c.Request.URL.RawQuery
			}
			span.SetAttributes(attribute.String("http.path", path))
			span.SetAttributes(attribute.String("http.method", c.Request.Method))
			span.SetAttributes(attribute.String("http.proto", c.Request.Proto))
			span.SetAttributes(attribute.String("ip", c.ClientIP()))
			if ipCountry := c.GetHeader("CF-IPCountry"); ipCountry != "" {
				span.SetAttributes(attribute.String("ipCountry", ipCountry))
			}
			span.SetAttributes(attribute.String("http.userAgent", c.Request.UserAgent()))
			span.SetAttributes(attribute.Int("http.req.body.size", int(c.Request.ContentLength)))
			c.Request = c.Request.WithContext(ctx)
			c.Next()
			span.SetAttributes(attribute.Int("http.resp.body.size", c.Writer.Size()))
			span.SetAttributes(attribute.Int("http.status_code", c.Writer.Status()))
		} else {
			c.Next()
		}
	}
}
