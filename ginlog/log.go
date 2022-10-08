package ginlog

import (
	"github.com/gin-gonic/gin"
	"github.com/zzzzer91/zlog"
	"strings"
	"time"
)

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		var err error
		if len(c.Errors) > 0 {
			err = c.Errors[len(c.Errors)-1].Err
		}
		ctx := c.Request.Context()
		msg := buildMsg(c, start)
		if err != nil {
			zlog.Ctx(ctx).WithError(err).Error(msg)
		} else {
			zlog.Ctx(ctx).Info(msg)
		}
	}
}

func buildMsg(c *gin.Context, start time.Time) string {
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
	sb.WriteString(time.Since(start).String())
	return sb.String()
}
