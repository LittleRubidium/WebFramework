package config

import (
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"path/filepath"
)

type HadeConfigProvider struct {
}

func (con *HadeConfigProvider) Register(c framework.Container) framework.NewInstance {
	return NewHadeConfig
}

func (con *HadeConfigProvider) Boot(c framework.Container) error {
	return nil
}

func (con *HadeConfigProvider) IsDefer() bool {
	return false
}

func (con *HadeConfigProvider) Params(c framework.Container) []interface{} {
	appService := c.MustMake(contract.AppKey).(contract.App)
	envService := c.MustMake(contract.EnvKey).(contract.Env)
	env := envService.AppEnv()
	configFolder := appService.ConfigFolder()
	envFolder := filepath.Join(configFolder, env)
	return []interface{}{c, envFolder, envService.All()}
}

func (con *HadeConfigProvider) Name() string {
	return contract.ConfigKey
}
