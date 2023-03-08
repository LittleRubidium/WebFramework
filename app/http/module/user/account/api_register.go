package account

import (
	"fmt"
	"github.com/gohade/hade/app/provider/user/account"
	"github.com/gohade/hade/framework/gin"
)

type registerParam struct {
	Username string `form:"username" binding:"required" json:"username"`
	Password string `form:"password" binding:"required" json:"password"`
	ConfirmedPassword string `form:"confirmedPassword" binding:"required" json:"confirmedPassword"`
}

func (api *AccountApi) Register(c *gin.Context) {
	userService := c.MustMake(account.UserKey).(account.Service)
	//logger := c.MustMake(contract.LogKey).(contract.Log)
	var param registerParam
	resp := map[string]string{}
	if err := c.ShouldBind(&param); err != nil {
		resp["error_message"] = "参数错误"
		c.ISetStatus(400).IJson(resp)
		//logger.Error(c, err.Error(),map[string]interface{}{
		//	"stack": fmt.Printf("%+v", err),
		//})
		c.ISetStatus(500).IText(err.Error())
		return
	}
	fmt.Println(param)
	resp = userService.Register(param.Username, param.Password, param.ConfirmedPassword)
	c.ISetOkStatus().IJson(resp)
}