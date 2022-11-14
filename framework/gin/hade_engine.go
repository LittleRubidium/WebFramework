package gin

import "github.com/gohade/hade/framework"

func (engine *Engine) SetContainer(c framework.Container) {
	engine.container = c
}

//engine实现container封装
func (engine *Engine) Bind(provider framework.ServiceProvider) error {
	return engine.container.Bind(provider)
}

//IsBand关键字凭证是否已经绑定服务提供者
func (engine *Engine) IsBind(key string) bool {
	return engine.container.IsBind(key)
}
