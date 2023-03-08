package bot

import (
	"github.com/gohade/hade/app/http/utils"
	"github.com/gohade/hade/app/provider/user/bot"
	"github.com/gohade/hade/framework/gin"
)

type updateParam struct {
	BotId       int    `form:"bot_id"`
	Title       string `form:"title"`
	Description string `form:"description"`
	Content     string `form:"content"`
}

func (ba *BotApi) Update(c *gin.Context) {
	resp := map[string]string{}
	botService := c.MustMake(bot.BotKey).(bot.Service)
	param := &updateParam{}
	if err := c.ShouldBind(param); err != nil {
		resp["error_message"] = "参数错误"
		c.ISetStatus(400).IJson(resp)
		return
	}
	data := map[string]interface{}{
		"bot_id":      param.BotId,
		"title":       param.Title,
		"description": param.Description,
		"content":     param.Content,
		"user_id":     utils.GetUser(c).Id,
	}
	resp = botService.Update(data)
	c.ISetOkStatus().IJson(resp)
}
