package util

import "strconv"

type Player struct {
	Id int
	BotId int
	BotCode string
	Sx,Sy int
	Steps []int
}

func (player *Player) checkTailIncrease(step int) bool {
	if step <= 10 {
		return true
	}
	return step % 3 == 1
}

func (player *Player) GetCells() []Cell {
	var res []Cell
	dx, dy := []int{-1,0,1,0},[]int{0,1,0,-1}
	x,y := player.Sx,player.Sy
	step := 0
	res = append(res,Cell{x,y})
	for _, d := range player.Steps {
		x += dx[d]
		y += dy[d]
		res = append(res,Cell{x,y})
		step++
		if !player.checkTailIncrease(step) {
			res = res[:len(res) - 1]
		}
	}
	return res
}

func (player *Player) GetStepString() string {
	res := ""
	for _,d := range player.Steps {
		res += strconv.Itoa(d)
	}
	return res
}