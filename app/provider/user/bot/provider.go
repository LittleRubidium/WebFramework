package bot

import "github.com/gohade/hade/framework"

type BotProvider struct {
}

func (bot *BotProvider) Register(c framework.Container) framework.NewInstance {
	return NewBotService
}

func (bot *BotProvider) Boot(c framework.Container) error {
	return nil
}

func (bot *BotProvider) Name() string {
	return BotKey
}

func (bot *BotProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

func (bot *BotProvider) IsDefer() bool {
	return true
}
