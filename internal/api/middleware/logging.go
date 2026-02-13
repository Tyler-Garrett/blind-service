package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		log.Printf("[%s] %s %s - %d - %v - %s",
			method,
			path,
			query,
			statusCode,
			latency,
			clientIP,
		)
	}
}
