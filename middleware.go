package gin_middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zzzzer91/gin-middleware/ginlogger"
	"github.com/zzzzer91/gin-middleware/ginmetrics"
	"github.com/zzzzer91/gin-middleware/ginrecovery"
	"github.com/zzzzer91/gin-middleware/gintrace"
)

func Default() []gin.HandlerFunc {
	return []gin.HandlerFunc{ginrecovery.Recovery(), ginlogger.Logger(), ginmetrics.Metrics(), gintrace.Trace()}
}
