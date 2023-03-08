package pk

import "github.com/gohade/hade/framework"

type PkProvider struct {
	framework.ServiceProvider
	container framework.Container
}

func (sp *PkProvider) Name() string {
	return PKKey
}
func (sp *PkProvider) Register(c framework.Container) framework.NewInstance {
	return NewPkService
}
func (sp *PkProvider) IsDefer() bool {
	return false
}
func (sp *PkProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}
func (sp *PkProvider) Boot(c framework.Container) error {
	return nil
}
