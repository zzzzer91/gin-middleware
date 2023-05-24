package ginauth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckingAuthorizationBearer(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := c.GetHeader("Authorization")
		fields := strings.Fields(s)
		if len(fields) != 2 || fields[0] != "Bearer" || fields[1] != token {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
