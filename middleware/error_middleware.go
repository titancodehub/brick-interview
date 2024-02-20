package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandler(errMap map[error]int) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, e := range c.Errors {
			err := e.Err
			val, ok := errMap[err]
			if !ok {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
				return
			}

			c.AbortWithStatusJSON(val, gin.H{
				"error": err.Error(),
			})
		}

		c.Next()
	}
}
