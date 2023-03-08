package consumer

import (
	"fmt"
	"github.com/gohade/hade/app/provider/user/account"
	"github.com/gohade/hade/app/provider/user/bot"
	"github.com/gohade/hade/app/utils/jwt"
	"github.com/gohade/hade/app/utils/restTemplate"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"github.com/gohade/hade/framework/gin"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	addPlayerUrl    = "http://127.0.0.1:3001/player/add/"
	removePlayerUrl = "http://127.0.0.1:3001/player/remove/"
)

type WebSocket struct {
	Upgrader *websocket.Upgrader
	Users    *sync.Map
	framework.Container
}

type Connect struct {
	Conn *websocket.Conn
	User *account.User
	Web  *WebSocket
	Game *Game
}

func (web *WebSocket) CreateConn(c *gin.Context) {
	fmt.Println("create conn")
	w, r := c.Writer, c.Request
	token := c.Param("token")
	userId := jwt.GetUserIdFromToken(token)
	userDB := &account.User{}
	tUser, ok := c.Get(strconv.Itoa(userId))
	userDB = tUser.(*account.User)
	if !ok {
		ormService := c.MustMake(contract.ORMKey).(contract.ORMService)
		db, err := ormService.GetDB()
		if err != nil {
			return
		}
		if db.Where("id=?", userId).First(userDB).Error == gorm.ErrRecordNotFound {
			return
		}
	}
	conn, err := web.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}
	con := &Connect{Conn: conn, User: userDB, Web: web}
	web.Users.Store(userId, con)
	go con.OnMessage()
}

func (web *WebSocket) StartGame(aId, aBotId, bId, bBotId int) {
	fmt.Println("start...")
	ormService := web.MustMake(contract.ORMKey).(contract.ORMService)
	db, err := ormService.GetDB()
	if err != nil {
		panic(err)
	}
	a, b := &account.User{}, &account.User{}
	var botA, botB *bot.Bot
	db.Where("id=?", aId).First(a)
	db.Where("id=?", bId).First(b)
	db.Table("bot").Where("id=?", aBotId).First(botA)
	db.Table("bot").Where("id=?", bBotId).First(botB)
	game := NewGame(13, 14, 20, aId, botA, bId, botB, web)
	game.CreateGameMap()
	if connA, ok := web.Users.Load(aId); ok {
		connA.(*Connect).Game = game
	}
	if connB, ok := web.Users.Load(bId); ok {
		connB.(*Connect).Game = game
	}
	respGame := map[string]interface{}{
		"a_id":    aId,
		"a_sx":    game.PlayerA.Sx,
		"a_sy":    game.PlayerA.Sy,
		"b_id":    bId,
		"b_sx":    game.PlayerB.Sx,
		"b_sy":    game.PlayerB.Sy,
		"gamemap": game.GetMap(),
	}
	respA := map[string]interface{}{
		"event":             "start-matching",
		"opponent_username": b.Username,
		"opponent_photo":    b.Photo,
		"game":              respGame,
	}
	if connA, ok := web.Users.Load(aId); ok {
		connA.(*Connect).SendMessage(respA)
	}
	respB := map[string]interface{}{
		"event":             "start-matching",
		"opponent_username": a.Username,
		"opponent_photo":    a.Photo,
		"game":              respGame,
	}
	fmt.Println(aBotId, " ", bBotId)
	if connB, ok := web.Users.Load(bId); ok {
		connB.(*Connect).SendMessage(respB)
	}
	time.Sleep(2 * time.Second)
	go game.Play()
}

func (conn *Connect) OnMessage() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
		conn.OnClose()
	}()
	for {
		_, message, err := conn.Conn.ReadMessage()
		if err != nil {
			return
		}
		fmt.Println("receive message!")
		data := gjson.Parse(string(message))
		event := data.Get("event").String()
		if strings.Compare(event, "start-matching") == 0 {
			conn.startMatching(int(data.Get("bot_id").Int()))
		} else if strings.Compare(event, "stop-matching") == 0 {
			conn.stopMatching()
		} else if strings.Compare(event, "move") == 0 {
			conn.move(int(data.Get("direction").Int()))
		}
	}
}

func (conn *Connect) OnClose() {
	fmt.Println("Close!")
	if conn.User != nil {
		conn.Web.Users.Delete(conn.User.Id)
	}
}

func (conn *Connect) startMatching(botId int) {
	data := url.Values{}
	data.Set("user_id", strconv.Itoa(conn.User.Id))
	data.Set("rating", strconv.Itoa(conn.User.Rating))
	data.Set("bot_id", strconv.Itoa(botId))
	restTemplate.PostForObject(addPlayerUrl, data)
}

func (conn *Connect) stopMatching() {
	data := url.Values{}
	data.Set("user_id", strconv.Itoa(conn.User.Id))
	restTemplate.PostForObject(removePlayerUrl, data)
}
func (conn *Connect) move(direction int) {
	if conn.Game.PlayerA.Id == conn.User.Id {
		if conn.Game.PlayerA.BotId == -1 {
			conn.Game.SetNextStepA(direction)
		}
	} else if conn.Game.PlayerB.Id == conn.User.Id {
		if conn.Game.PlayerB.BotId == -1 {
			conn.Game.SetNextStepB(direction)
		}
	}
}

func (conn *Connect) SendMessage(message interface{}) {
	conn.Conn.WriteJSON(message)
}
