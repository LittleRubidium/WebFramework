package pk

import (
	"fmt"
	"github.com/gohade/hade/app/http/consumer"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
)

type PkService struct {
	container framework.Container
	logger contract.Log
	configer contract.Config
}

func (pk *PkService) StartGame(aId, aBotId, bId, bBotId int) string {
	fmt.Println("start game ",aId," ",bId)
	consumer.Web.StartGame(aId,aBotId,bId,bBotId)
	return "start game success"
}

func (pk *PkService) ReceiveBotMove(userId, direction int) string {
	fmt.Println("receive bot move: ",userId, " ", direction)
	web := consumer.Web
	if conn,ok := web.Users.Load(userId); ok {
		game := conn.(*consumer.Connect).Game
		if game != nil {
			if game.PlayerA.Id == userId {
				game.SetNextStepA(direction)
			}else if game.PlayerB.Id == userId {
				game.SetNextStepB(direction)
			}
		}
	}
	return "receive bot move success"
}

func NewPkService(params []interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	logger := container.MustMake(contract.LogKey).(contract.Log)
	configer := container.MustMake(contract.ConfigKey).(contract.Config)
	return &PkService{container: container, logger: logger,configer: configer}, nil
}
