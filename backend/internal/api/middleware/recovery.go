package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Recovery returns a middleware that recovers from panics and logs the stack trace.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				logrus.WithFields(logrus.Fields{
					"error":  r,
					"stack":  stack,
					"method": c.Request.Method,
					"path":   c.Request.URL.Path,
				}).Error("panic recovered")

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"ok":      false,
					"message": "internal server error",
				})
			}
		}()
		c.Next()
	}
}
