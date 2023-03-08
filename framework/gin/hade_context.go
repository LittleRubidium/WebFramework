package gin

import (
	"context"
	"github.com/gohade/hade/framework"
)

func (ctx *Context) BaseContext() context.Context {
	return ctx.Request.Context()
}
func (ctx *Context) GetContainer() framework.Container {
	return ctx.container
}
