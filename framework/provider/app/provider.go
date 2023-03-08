package app

import (
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
)

//HadeAppProvider
type HadeAppProvider struct {
	BaseFolder string
}

//注册HadeApp方法
func (ha *HadeAppProvider) Register(c framework.Container) framework.NewInstance {
	return NewHadeApp
}

//启动调用
func (ha *HadeAppProvider) Boot(c framework.Container) error {
	return nil
}

//是否延迟实例化
func (ha *HadeAppProvider) IsDefer() bool {
	return false
}

//获取初始化参数
func (ha *HadeAppProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c, ha.BaseFolder}
}

//获取字符串凭证
func (ha *HadeAppProvider) Name() string {
	return contract.AppKey
}
