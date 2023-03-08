package cache

import (
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"github.com/gohade/hade/framework/provider/cache/services"
	"strings"
)

type HadeCacheProvider struct {
	Driver string
}

func (cache *HadeCacheProvider) Register(c framework.Container) framework.NewInstance {
	if cache.Driver == "" {
		tcs, err := c.Make(contract.ConfigKey)
		if err != nil {
			return services.NewMemoryCache
		}
		cs := tcs.(contract.Config)
		cache.Driver = strings.ToLower(cs.GetString("cache.driver"))
	}

	switch cache.Driver {
	case "redis":
		return services.NewRedisCache
	case "memory":
		return services.NewMemoryCache
	default:
		return services.NewMemoryCache
	}
}

func (cache *HadeCacheProvider) Boot(c framework.Container) error {
	return nil
}

func (cache *HadeCacheProvider) IsDefer() bool {
	return true
}

func (cache *HadeCacheProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

func (cache *HadeCacheProvider) Name() string {
	return contract.CacheKey
}
