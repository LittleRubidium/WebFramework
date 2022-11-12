package kernel

import (
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"github.com/gohade/hade/framework/gin"
)

//提供web引擎
type HadeKernelProvider struct {
	HttpEngine *gin.Engine
}

func (hk *HadeKernelProvider) Register(c framework.Container) framework.NewInstance {
	return NewHadeKernelProvider
}

func (hk *HadeKernelProvider) Boot(c framework.Container) error {
	if hk.HttpEngine == nil {
		hk.HttpEngine = gin.Default()
	}
	hk.HttpEngine.SetContainer(c)
	return nil
}

func (hk *HadeKernelProvider) IsDefer() bool {
	return false
}

func (hk *HadeKernelProvider) Params(c framework.Container) []interface{} {
	return []interface{}{hk.HttpEngine}
}

func (hk *HadeKernelProvider) Name() string {
	return contract.KernelKey
}
