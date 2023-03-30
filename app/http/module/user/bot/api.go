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
	api := r.Group("/api")
	{
		user := api.Group("/user")
		{
			rbot := user.Group("/bot")
			{
				rbot.POST("/", bot.Add)
				rbot.GET("/", bot.GetList)
				rbot.DELETE("/", bot.Remove)
				rbot.PUT("/", bot.Update)
			}
		}
	}
	return nil
}
