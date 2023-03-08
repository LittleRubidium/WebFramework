package pk

const PKKey = "pk"

type Service interface {
	ReceiveBotMove(userId, direction int) string
	StartGame(aId,aBotId,bId,bBotId int) string
}
