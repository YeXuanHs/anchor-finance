package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger returns a middleware that logs request details using logrus.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		entry := logrus.WithFields(logrus.Fields{
			"status":  status,
			"method":  method,
			"path":    path,
			"query":   query,
			"ip":      clientIP,
			"latency": latency.String(),
			"length":  c.Writer.Size(),
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else if status >= 500 {
			entry.Error("server error")
		} else if status >= 400 {
			entry.Warn("client error")
		} else {
			entry.Info("request")
		}
	}
}
