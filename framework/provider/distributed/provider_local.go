package distributed

import (
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
)

type LocalDistributedProvider struct {
}

func (l *LocalDistributedProvider) Boot(c framework.Container) error {
	return nil
}

func (l *LocalDistributedProvider) Register(c framework.Container) framework.NewInstance {
	return NewLocalDistributedService
}

func (l *LocalDistributedProvider) IsDefer() bool {
	return false
}

func (l *LocalDistributedProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

func (l *LocalDistributedProvider) Name() string {
	return contract.DistributedKey
}
