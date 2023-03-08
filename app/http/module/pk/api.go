package pk

import (
	"github.com/gohade/hade/app/provider/pk"
	"github.com/gohade/hade/framework/gin"
)

type PkApi struct {
}

func Register(r *gin.Engine) {
	pkApi := &PkApi{}
	r.Bind(&pk.PkProvider{})
	r.POST("/pk/start/game/",pkApi.StartGame)
	r.POST("/pk/receive/bot/move/",pkApi.ReceiveBotMove)
}
