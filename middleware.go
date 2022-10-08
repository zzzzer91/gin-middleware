package gin_middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zzzzer91/gin-middleware/ginlog"
	"github.com/zzzzer91/gin-middleware/ginmetric"
	"github.com/zzzzer91/gin-middleware/ginrecovery"
	"github.com/zzzzer91/gin-middleware/gintrace"
)

func Default() []gin.HandlerFunc {
	return []gin.HandlerFunc{ginrecovery.Recovery(), ginlog.Log(), ginmetric.Metrics(), gintrace.Trace()}
}
