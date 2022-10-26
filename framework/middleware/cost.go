package middleware

import (
	"github.com/gohade/hade/framework/gin"
	"log"
	"time"
)

func Cost() gin.HandlerFunc {
	return func(c *gin.Context) {
		//记录开始时间
		start := time.Now()

		//使用next执行具体业务逻辑
		c.Next()

		//记录结束时间
		end := time.Now()
		cost := end.Sub(start)
		log.Printf("api uri: %v, cost: %v", c.Request.RequestURI, cost.Seconds())
	}
}
