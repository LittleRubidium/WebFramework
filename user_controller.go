package main

import (
	"github.com/gohade/hade/framework/gin"
	"time"
)

func UserLoginController(c *gin.Context) {
	foo, _ := c.DefaultQueryString("foo", "def")
	time.Sleep(10 * time.Second)
	c.ISetStatus(200).IJson("ok UserLoginController: " + foo)
}
