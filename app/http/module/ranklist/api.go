package ranklist

import (
	"github.com/gohade/hade/app/provider/ranklist"
	"github.com/gohade/hade/framework/gin"
)

type RankApi struct {
}

func Register(r *gin.Engine) error {
	rankApi := &RankApi{}
	r.Bind(&ranklist.RankListProvider{})
	r.GET("/api/ranklist/getlist/", rankApi.GetRankList)
	return nil
}
