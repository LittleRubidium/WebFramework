package bot

import (
	"github.com/gohade/hade/app/provider/user/bot"
	"github.com/gohade/hade/framework/gin"
)

type removeParam struct {
	BotId int `form:"bot_id"`
}

func (ba *BotApi) Remove(c *gin.Context) {
	param := &removeParam{}
	resp := map[string]string{}
	if err := c.ShouldBind(param); err != nil {
		resp["error_message"] = "参数错误"
		c.ISetStatus(400).IJson(resp)
		return
	}
	botService := c.MustMake(bot.BotKey).(bot.Service)
	resp = botService.Remove(param.BotId)
	c.ISetOkStatus().IJson(resp)
}
