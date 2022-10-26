package middleware

import "github.com/gohade/hade/framework/gin"

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.ISetStatus(500).IJson(err)
			}
		}()
		//使用next执行具体的业务逻辑
		c.Next()

	}
}
