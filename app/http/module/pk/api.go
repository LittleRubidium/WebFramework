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
	pk := r.Group("/pk")
	{
		pk.POST("/", pkApi.StartGame)
		pk.PUT("/", pkApi.ReceiveBotMove)
	}
}
