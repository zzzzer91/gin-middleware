package gintrace

import (
	"github.com/gin-gonic/gin"
	"github.com/zzzzer91/zlog"
	"go.opentelemetry.io/otel/attribute"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := zlog.StartTracing(c.Request.Context(), c.Request.Method+" "+c.FullPath())
		defer span.End()
		path := c.Request.URL.Path
		if len(c.Request.URL.RawQuery) > 0 {
			path += "?" + c.Request.URL.RawQuery
		}
		span.SetAttributes(attribute.String("path", path))
		span.SetAttributes(attribute.String("method", c.Request.Method))
		span.SetAttributes(attribute.String("ip", c.ClientIP()))
		span.SetAttributes(attribute.String("userAgent", c.Request.UserAgent()))
		span.SetAttributes(attribute.Int("req.body.size", int(c.Request.ContentLength)))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
		span.SetAttributes(attribute.Int("resp.body.size", c.Writer.Size()))
	}
}
