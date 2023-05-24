package gincors

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Fix: https://github.com/gin-gonic/gin/issues/2406
type corsWriter struct {
	gin.ResponseWriter
	allowOrigin string
	maxAge      string
}

func (w *corsWriter) WriteHeader(statusCode int) {
	if statusCode == http.StatusOK {
		w.Header().Add("Access-Control-Allow-Origin", w.allowOrigin)
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Add("Access-Control-Max-Age", w.maxAge)
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

// for single route
func Cors(allowOrigin string, maxAge time.Duration) gin.HandlerFunc {
	maxAgeStr := fmt.Sprintf("%.f", maxAge.Seconds())

	return func(c *gin.Context) {
		c.Writer = &corsWriter{ResponseWriter: c.Writer, allowOrigin: allowOrigin, maxAge: maxAgeStr}
		c.Next()
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
