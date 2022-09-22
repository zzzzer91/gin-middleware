package gintrace

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zzzzer91/gopkg/tracex"
	"github.com/zzzzer91/gopkg/zlog"
	"go.opentelemetry.io/otel/attribute"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracex.StartTracing(newCtx(c), c.Request.Method+":"+c.FullPath())
		defer span.End()
		span.SetAttributes(attribute.String("uri", c.Request.RequestURI))
		span.SetAttributes(attribute.String("method", c.Request.Method))
		span.SetAttributes(attribute.String("ip", c.ClientIP()))
		span.SetAttributes(attribute.String("userAgent", c.Request.UserAgent()))
		span.SetAttributes(attribute.Int("contentLength", int(c.Request.ContentLength)))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func newCtx(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	if requestID := c.GetHeader(httpHeaderFieldNameRequestID); requestID != "" {
		ctx = context.WithValue(ctx, zlog.EntityFieldNameRequestID, requestID)
	}
	return ctx
}
