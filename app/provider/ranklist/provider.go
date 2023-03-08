package ranklist

import "github.com/gohade/hade/framework"

type RankListProvider struct {
	framework.ServiceProvider
	c framework.Container
}
func (sp *RankListProvider) Name() string {
	return RanklistKey
}
func (sp *RankListProvider) Register(c framework.Container) framework.NewInstance {
	return NewRankListService
}
func (sp *RankListProvider) IsDefer() bool {
	return false
}
func (sp *RankListProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}
func (sp *RankListProvider) Boot(c framework.Container) error {
	return nil
}
