package pk

import (
	"github.com/gohade/hade/app/provider/pk"
	"github.com/gohade/hade/framework/gin"
	"strconv"
)

func (pa *PkApi) StartGame(c *gin.Context) {

	aId, _ := strconv.Atoi(c.PostForm("a_id"))
	aBotId, _ := strconv.Atoi(c.PostForm("a_bot_id"))
	bId, _ := strconv.Atoi(c.PostForm("b_id"))
	bBotId, _ := strconv.Atoi(c.PostForm("b_bot_id"))
	PkService := c.MustMake(pk.PKKey).(pk.Service)
	resp := PkService.StartGame(aId, aBotId, bId, bBotId)
	c.ISetOkStatus().IText(resp)
}
