package account

import (
	"github.com/gohade/hade/app/http/utils"
	"github.com/gohade/hade/framework/gin"
	"strconv"
)

func (api *AccountApi) Info(c *gin.Context) {
	resp := map[string]string{}
	user := utils.GetUser(c)
	resp["error_message"] = "success"
	resp["id"] = strconv.Itoa(user.Id)
	resp["username"] = user.Username
	resp["photo"] = user.Photo
	c.ISetOkStatus().IJson(resp)
}