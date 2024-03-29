package ginlog

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zzzzer91/zlog"
)

const (
	httpHeaderFieldNameRequestID = "X-Request-ID"
)

func Log(isLogInfo bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx := newCtx(c)
		basicMsg := buildBasicMsg(c)
		if isLogInfo {
			zlog.Ctx(ctx).Info(buildBeginMsg(basicMsg))
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
		used := time.Since(start)
		ctx = c.Request.Context()
		var err error
		if len(c.Errors) > 0 {
			err = c.Errors[len(c.Errors)-1].Err
		}
		if err != nil {
			zlog.Ctx(ctx).WithError(err).Error(buildErrorMsg(basicMsg, used))
		} else {
			if isLogInfo {
				zlog.Ctx(ctx).Info(buildCostMsg(basicMsg, used))
			}
		}
	}
}

func newCtx(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	if requestID := c.GetHeader(httpHeaderFieldNameRequestID); requestID != "" {
		ctx = context.WithValue(ctx, zlog.EntityFieldNameRequestID, requestID)
	}
	ctx = zlog.AddLogIDToCtx(ctx)
	return ctx
}

func buildBasicMsg(c *gin.Context) string {
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
	sb.WriteByte('`')
	return sb.String()
}

func buildBeginMsg(basicMsg string) string {
	return basicMsg + ", started"
}

func buildCostMsg(basicMsg string, used time.Duration) string {
	return basicMsg + ", used " + used.String()
}

func buildErrorMsg(basicMsg string, used time.Duration) string {
	return basicMsg + ", used " + used.String() + ", error!"
}
