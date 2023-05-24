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
			c.Header("Access-Control-Allow-Origin", allowOrigin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
			c.Header("Access-Control-Max-Age", fmt.Sprintf("%.f", maxAge.Seconds()))
		}
	}
}

// for global usage
func CorsGlobal(allowOrigin string, maxAge time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			c.Header("Access-Control-Allow-Origin", allowOrigin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE")
			c.Header("Access-Control-Max-Age", fmt.Sprintf("%.f", maxAge.Seconds()))
			return
		}
		c.Next()
	}
}
