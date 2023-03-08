package account

import (
	"github.com/gohade/hade/app/provider/user/account"
	"github.com/gohade/hade/framework/gin"
)

type loginParams struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (api *AccountApi) Login(c *gin.Context) {
	param := &loginParams{}
	resp := map[string]string{}
	if err := c.ShouldBind(param); err != nil {
		resp["error_message"] = "参数错误"
		c.ISetStatus(400).IJson(resp)
	}
	userService := c.MustMake(account.UserKey).(account.Service)
	resp = userService.Login(param.Username, param.Password)
	c.ISetOkStatus().IJson(resp)
}
