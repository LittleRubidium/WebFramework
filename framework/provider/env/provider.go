package env

import (
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
)

type HadeEnvProvider struct {
	Folder string
}

func (e *HadeEnvProvider) Register(c framework.Container) framework.NewInstance {
	return NewHadeEnv
}

func (e *HadeEnvProvider) Boot(c framework.Container) error {
	app := c.MustMake(contract.AppKey).(contract.App)
	e.Folder = app.BaseFolder()
	return nil
}

func (e *HadeEnvProvider) IsDefer() bool {
	return false
}

func (e *HadeEnvProvider) Params(c framework.Container) []interface{} {
	return []interface{}{e.Folder}
}

func (e *HadeEnvProvider) Name() string {
	return contract.EnvKey
}
