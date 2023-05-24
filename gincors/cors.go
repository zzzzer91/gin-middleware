package gincors

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// for single route
func Cors(allowOrigin string, maxAge time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() == http.StatusOK {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
			c.Writer.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%.f", maxAge.Seconds()))
		}
	}
}

// for global usage
func CorsGlobal(allowOrigin string, maxAge time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%.f", maxAge.Seconds()))
			return
		}
		c.Next()
	}
}
