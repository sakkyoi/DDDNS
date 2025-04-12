package middleware

import (
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop Timer
		end := time.Now()
		latency := end.Sub(start)

		// information from the request
		clientIp := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		// use LogFormatterParams to get colors for status code and method
		params := gin.LogFormatterParams{
			Method:     method,
			StatusCode: statusCode,
		}

		// if the request has a query string, append it to the path
		if raw != "" {
			path += "?" + raw
		}

		// colors
		statusColor := params.StatusCodeColor()
		methodColor := params.MethodColor()
		resetColor := params.ResetColor()

		log.Debugf("[API] |%s %3d %s| %13v | %15s |%s %-7s %s %#v",
			statusColor, statusCode, resetColor,
			latency,
			clientIp,
			methodColor, method, resetColor,
			path,
		)
	}
}
