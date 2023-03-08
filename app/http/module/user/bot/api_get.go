package bot

import (
	"github.com/gohade/hade/app/http/utils"
	"github.com/gohade/hade/app/provider/user/bot"
	"github.com/gohade/hade/framework/gin"
)

func (ba *BotApi) GetList(c *gin.Context) {
	botService := c.MustMake(bot.BotKey).(bot.Service)
	user := utils.GetUser(c)
	resp := botService.GetList(user.Id)
	c.ISetOkStatus().IJson(resp)
}
