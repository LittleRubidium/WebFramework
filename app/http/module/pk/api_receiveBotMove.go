package pk

import (
	"github.com/gohade/hade/app/provider/pk"
	"github.com/gohade/hade/framework/gin"
)

type botMoveParam struct {
	userId int `form:"user_id"`
	direction int `form:"direction"`
}

func (pa *PkApi) ReceiveBotMove(c *gin.Context) {
	param := &botMoveParam{}
	if err := c.ShouldBind(param); err != nil {
		c.ISetStatus(400).IText("参数错误")
	}
	pkService := c.MustMake(pk.PKKey).(pk.Service)
	resp := pkService.ReceiveBotMove(param.userId,param.direction)
	c.ISetOkStatus().IText(resp)
}
