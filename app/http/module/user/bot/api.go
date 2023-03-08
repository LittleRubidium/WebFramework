package bot

import (
	"github.com/gohade/hade/app/provider/user/bot"
	"github.com/gohade/hade/framework/gin"
)

type BotApi struct {
}

func Register(r *gin.Engine) error {
	r.Bind(&bot.BotProvider{})
	bot := &BotApi{}
	r.POST("/api/user/bot/add/",bot.Add)
	r.GET("/api/user/bot/getlist/",bot.GetList)
	r.POST("/api/user/bot/remove/",bot.Remove)
	r.POST("/api/user/bot/update/",bot.Update)
	return nil
}
