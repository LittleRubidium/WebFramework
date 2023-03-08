package account

import (
	"fmt"
	"github.com/gohade/hade/app/provider/user"
	"github.com/gohade/hade/framework/contract"
	"github.com/gohade/hade/framework/gin"
)

type registerParam struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

func (api *AccountApi) Register(c *gin.Context) {
	userService := c.MustMake(user.UserKey).(user.Service)
	logger := c.MustMake(contract.LogKey).(contract.Log)

	param := &registerParam{}
	if err := c.ShouldBindJSON(param); err != nil {
		c.ISetStatus(400).IText("参数错误")
		logger.Error(c, err.Error(),map[string]interface{}{
			"stack": fmt.Printf("%+v", err),
		})
		c.ISetStatus(500).IText(err.Error())
		return
	}
	resp := userService.Register(param.Username, param.Password, param.ConfirmPassword)
	c.ISetOkStatus().IJson(resp)
}