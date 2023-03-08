package bot

import (
	"fmt"
	"github.com/gohade/hade/app/http/utils"
	"github.com/gohade/hade/app/provider/user/bot"
	"github.com/gohade/hade/framework/gin"
)

type addBotParams struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	Content     string `form:"content"`
}

func (ba *BotApi) Add(c *gin.Context) {
	botService := c.MustMake(bot.BotKey).(bot.Service)
	param := &addBotParams{}
	resp := map[string]string{}
	if err := c.ShouldBind(param); err != nil {
		resp["error_message"] = "参数错误"
		c.ISetStatus(400).IJson(resp)
		return
	}
	fmt.Println(param)
	user := utils.GetUser(c)
	addParams := map[string]interface{}{
		"title":       param.Title,
		"description": param.Description,
		"content":     param.Content,
		"userId":      user.Id,
	}
	resp = botService.Add(addParams)
	c.ISetOkStatus().IJson(resp)
}
