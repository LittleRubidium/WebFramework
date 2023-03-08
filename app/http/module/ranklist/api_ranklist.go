package ranklist

import (
	"github.com/gohade/hade/app/provider/ranklist"
	"github.com/gohade/hade/framework/gin"
	"strconv"
)

func (r *RankApi) GetRankList(c *gin.Context) {
	rankService := c.MustMake(ranklist.RanklistKey).(ranklist.Service)
	page, _ := strconv.Atoi(c.Query("page"))
	resp := rankService.GetRankList(page)
	c.ISetOkStatus().IJson(resp)
}
