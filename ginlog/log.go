package ginlog

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zzzzer91/zlog"
)

const (
	httpHeaderFieldNameRequestID = "X-Request-ID"
)

func Log(isLogInfo bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Request = c.Request.WithContext(newCtx(c))
		c.Next()
		used := time.Since(start)
		ctx := c.Request.Context()
		var err error
		if len(c.Errors) > 0 {
			err = c.Errors[len(c.Errors)-1].Err
		}
		if err != nil {
			zlog.Ctx(ctx).WithError(err).Error(buildErrorMsg(c, used))
		} else {
			if isLogInfo {
				zlog.Ctx(ctx).Info(buildInfoMsg(c, used))
			}
		}
	}
}

func newCtx(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	if requestID := c.GetHeader(httpHeaderFieldNameRequestID); requestID != "" {
		ctx = context.WithValue(ctx, zlog.EntityFieldNameRequestID, requestID)
	} else {
		ctx = context.WithValue(ctx, zlog.EntityFieldNameRequestID, uuid.New().String())
	}
	return ctx
}

func buildInfoMsg(c *gin.Context, used time.Duration) string {
	sb := strings.Builder{}
	sb.WriteString(c.ClientIP())
	sb.WriteByte(' ')
	sb.WriteString(c.Request.Method)
	sb.WriteString(" `")
	sb.WriteString(c.Request.URL.Path)
	if len(c.Request.URL.RawQuery) > 0 {
		sb.WriteByte('?')
		sb.WriteString(c.Request.URL.RawQuery)
	}
	sb.WriteString("`, used ")
	sb.WriteString(used.String())
	return sb.String()
}

func buildErrorMsg(c *gin.Context, used time.Duration) string {
	sb := strings.Builder{}
	sb.WriteString(c.ClientIP())
	sb.WriteByte(' ')
	sb.WriteString(c.Request.Method)
	sb.WriteString(" `")
	sb.WriteString(c.Request.URL.Path)
	if len(c.Request.URL.RawQuery) > 0 {
		sb.WriteByte('?')
		sb.WriteString(c.Request.URL.RawQuery)
	}
	sb.WriteString("`, used ")
	sb.WriteString(used.String())
	sb.WriteString(", error!")
	return sb.String()
}
