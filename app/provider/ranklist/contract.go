package ranklist

const RanklistKey = "ranklist"

type Service interface {
	GetRankList(int) map[string]interface{}
}
