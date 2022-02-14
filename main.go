package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/random"
	"time"
)

const (
	//CorrelationID is a request id unique to the request being made
	CorrelationID = "X-Correlation-ID"
)

func main() {
	r := gin.New()
	r.Use(addCorrelationID())
	r.Use(logMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	recordR := r.Group("/record")
	{
		recordR.POST("")
	}

	r.Run(":8080")
}

func logMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 自定义日志输出格式
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC3339Nano),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

func addCorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// generate correlation id
		id := c.Request.Header.Get(CorrelationID)
		var newID string
		if id == "" {
			//generate a random number
			newID = random.String(12)
		} else {
			newID = id
		}
		c.Request.Header.Set(CorrelationID, newID)
		c.Header(CorrelationID, newID)
		c.Next()
	}
}
