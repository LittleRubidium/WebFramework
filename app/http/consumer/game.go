package consumer

import (
	"fmt"
	"github.com/gohade/hade/app/provider/record"
	"github.com/gohade/hade/app/provider/user/account"
	"github.com/gohade/hade/app/provider/user/bot"
	"github.com/gohade/hade/app/utils/restTemplate"
	"github.com/gohade/hade/framework/contract"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	addBotUrl = "http://127.0.0.1:3002/bot/add/"
)

var (
	dx = [4]int{-1, 0, 1, 0}
	dy = [4]int{0, 1, 0, -1}
)

type Game struct {
	rows, cols           int
	PlayerA, PlayerB     *Player
	nextStepA, nextStepB int
	innerWallsCount      int
	g                    [][]int
	lock                 sync.RWMutex
	loser                string
	web                  *WebSocket
	status               string
}

func NewGame(rows, cols, innerWallsCount int, idA int, botA *bot.Bot, idB int, botB *bot.Bot, web *WebSocket) *Game {
	botIdA, botIdB := -1, -1
	botCodeA, botCodeB := "", ""
	if botA != nil {
		botIdA = botA.Id
		botCodeA = botA.Content
	}
	if botB != nil {
		botIdB = botB.Id
		botCodeB = botB.Content
	}
	playerA := &Player{
		Id:      idA,
		BotId:   botIdA,
		BotCode: botCodeA,
		Sx:      rows - 2,
		Sy:      1,
		Steps:   []int{},
	}
	playerB := &Player{
		Id:      idB,
		BotId:   botIdB,
		BotCode: botCodeB,
		Sx:      1,
		Sy:      cols - 2,
		Steps:   []int{},
	}
	g := make([][]int, rows)
	for i := 0; i < rows; i++ {
		g[i] = make([]int, cols)
	}
	return &Game{
		rows:            rows,
		cols:            cols,
		PlayerA:         playerA,
		PlayerB:         playerB,
		nextStepA:       -1,
		nextStepB:       -1,
		innerWallsCount: innerWallsCount,
		lock:            sync.RWMutex{},
		g:               g,
		web:             web,
		status:          "playing",
	}
}

func (g *Game) checkConnectivity(sx, sy, tx, ty int) bool {
	if sx == tx && sy == ty {
		return true
	}
	g.g[sx][sy] = 1
	for i := 0; i < 4; i++ {
		x, y := sx+dx[i], sy+dy[i]
		if x >= 0 && x < g.rows && y >= 0 && y < g.cols && g.g[x][y] == 0 {
			if g.checkConnectivity(x, y, tx, ty) {
				g.g[sx][sy] = 0
				return true
			}
		}
	}
	g.g[sx][sy] = 0
	return false
}

func (g *Game) drawMap() bool {
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			g.g[i][j] = 0
		}
	}
	for r := 0; r < g.rows; r++ {
		g.g[r][0] = 1
		g.g[r][g.cols-1] = 1
	}
	for c := 0; c < g.cols; c++ {
		g.g[0][c] = 1
		g.g[g.rows-1][c] = 1
	}
	for i := 0; i < g.innerWallsCount/2; i++ {
		for j := 0; j < 1000; j++ {
			r, c := rand.Intn(g.rows), rand.Intn(g.cols)
			if g.g[r][c] == 1 || g.g[g.rows-1-r][g.cols-1-c] == 1 {
				continue
			}
			if (r == g.rows-2 && c == 1) || (r == 1 && c == g.cols-2) {
				continue
			}
			g.g[r][c] = 1
			g.g[g.rows-1-r][g.cols-1-c] = 1
			break
		}
	}
	return g.checkConnectivity(g.rows-2, 1, 1, g.cols-2)
}

func (g *Game) CreateGameMap() {
	for i := 0; i < 1000; i++ {
		if g.drawMap() {
			break
		}
	}
}

func (g *Game) SetNextStepA(nextStepA int) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.nextStepA = nextStepA
}

func (g *Game) SetNextStepB(nextStepB int) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.nextStepB = nextStepB
}

func (g *Game) nextStep() bool {
	time.Sleep(time.Millisecond * 100)
	g.sendBotCode(g.PlayerA)
	g.sendBotCode(g.PlayerB)
	for i := 0; i < 50; i++ {
		time.Sleep(time.Millisecond * 100)
		g.lock.RLock()
		if g.nextStepA != -1 && g.nextStepB != -1 {
			g.PlayerA.Steps = append(g.PlayerA.Steps, g.nextStepA)
			g.PlayerB.Steps = append(g.PlayerB.Steps, g.nextStepB)
			g.lock.RUnlock()
			return true
		}
		g.lock.RUnlock()
	}
	return false
}

func (g *Game) sendBotCode(player *Player) {
	if player.BotId == -1 {
		return
	}
	data := url.Values{}
	data.Set("user_id", strconv.Itoa(player.Id))
	data.Set("bot_code", player.BotCode)
	data.Set("input", g.getInput(player))
	err := restTemplate.PostForObject(addBotUrl, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (g *Game) getInput(player *Player) string {
	me, you := &Player{}, &Player{}
	if g.PlayerA.Id == player.Id {
		me = g.PlayerA
		you = g.PlayerB
	} else {
		me = g.PlayerB
		you = g.PlayerA
	}
	return g.getMapString() +
		"#" + strconv.Itoa(me.Sx) +
		"#" + strconv.Itoa(me.Sy) +
		"#(" + me.GetStepString() + ")#" +
		strconv.Itoa(you.Sx) + "#" +
		strconv.Itoa(you.Sy) + "#(" +
		you.GetStepString() + ")"
}

func (g *Game) getMapString() string {
	res := ""
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			res += strconv.Itoa(g.g[i][j])
		}
	}
	return res
}

func (g *Game) sendAllMessage(resp interface{}) {
	if connA, ok := g.web.Users.Load(g.PlayerA.Id); ok {
		connA.(*Connect).Conn.WriteJSON(resp)
	}
	if connB, ok := g.web.Users.Load(g.PlayerB.Id); ok {
		connB.(*Connect).Conn.WriteJSON(resp)
	}
}

func (g *Game) sendMove() {
	g.lock.Lock()
	defer g.lock.Unlock()
	resp := map[string]interface{}{
		"event":       "move",
		"a_direction": g.nextStepA,
		"b_direction": g.nextStepB,
	}
	g.sendAllMessage(resp)
	g.nextStepA, g.nextStepB = -1, -1
}

func (g *Game) sendResult() {
	resp := map[string]string{
		"event": "result",
		"loser": g.loser,
	}
	g.saveToDatabase()
	g.sendAllMessage(resp)
}

func (g *Game) saveToDatabase() {
	ormService := g.web.MustMake(contract.ORMKey).(contract.ORMService)
	db, err := ormService.GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	a, b := &account.User{}, &account.User{}
	db.Where("id=?", g.PlayerA.Id).First(a)
	db.Where("id=?", g.PlayerB.Id).First(b)

	if strings.Compare("A", g.loser) == 0 {
		a.Rating -= 2
		b.Rating += 5
	} else if strings.Compare("B", g.loser) == 0 {
		a.Rating += 5
		b.Rating -= 2
	}
	db.Where("id=?", a.Id).Updates(a)
	db.Where("id=?", b.Id).Updates(b)
	newRecord := record.NewRecord(g.PlayerA.Id,
		g.PlayerA.Sx,
		g.PlayerA.Sy,
		g.PlayerB.Id,
		g.PlayerB.Sx,
		g.PlayerB.Sy,
		g.PlayerA.GetStepString(),
		g.PlayerB.GetStepString(),
		g.getMapString(),
		g.loser,
		time.Now())
	db.Table("record").Create(newRecord)
}

func (g *Game) checkValid(ca, cb []Cell) bool {
	n := len(ca)
	cell := ca[n-1]
	if g.g[cell.X][cell.Y] == 1 {
		return false
	}
	for i := 0; i < n-1; i++ {
		if ca[i].X == cell.X && ca[i].Y == cell.Y {
			return false
		}
	}
	for i := 0; i < n; i++ {
		if cb[i].X == cell.X && cb[i].Y == cell.Y {
			return false
		}
	}
	return true
}

func (g *Game) judge() {
	ca, cb := g.PlayerA.GetCells(), g.PlayerB.GetCells()
	va, vb := g.checkValid(ca, cb), g.checkValid(cb, ca)
	if !va || !vb {
		g.status = "finished"
		if !va && !vb {
			g.loser = "all"
		} else if !va {
			g.loser = "A"
		} else {
			g.loser = "B"
		}
	}
}

func (g *Game) Play() {
	for i := 0; i < 1000; i++ {
		if g.nextStep() {
			g.judge()
			if g.status == "playing" {
				g.sendMove()
			} else {
				g.sendResult()
				break
			}
		} else {
			g.status = "finished"
			g.lock.RLock()
			if g.nextStepA == -1 && g.nextStepB == -1 {
				g.loser = "all"
			} else if g.nextStepA == -1 {
				g.loser = "A"
			} else {
				g.loser = "B"
			}
			g.lock.RUnlock()
			g.sendResult()
			break
		}
	}
}

func (g *Game) GetMap() [][]int {
	return g.g
}
