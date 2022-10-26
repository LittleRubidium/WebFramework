package middleware

import (
	"context"
	"fmt"
	"github.com/gohade/hade/framework/gin"
	"log"
	"time"
)

func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		finish, pinicChan := make(chan struct{}, 1), make(chan interface{}, 1)
		//执行业务逻辑签预操作:初始化超时context
		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					pinicChan <- p
				}
			}()
			//使用next执行具体的业务逻辑
			c.Next()

			finish <- struct{}{}
		}()
		//执行业务逻辑后操作
		select {
		case p := <-pinicChan:
			c.ISetStatus(500).IJson("time out")
			log.Println(p)
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.ISetStatus(500).IJson("time out")
		}
	}
}
