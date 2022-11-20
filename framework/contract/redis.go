package contract

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gohade/hade/framework"
)

const RedisKey = "hade:redis"

type RedisOption func(container framework.Container, config *RedisConfig) error

//为hade定义的Redis配置结构
type RedisConfig struct {
	*redis.Options
}

//redis服务
type RedisService interface {
	GetClient(option ...RedisOption) (*redis.Client, error)
}

//用来唯一标识一个redisConfig配置
func (conf *RedisConfig) UniqKey() string {
	return fmt.Sprintf("%v_%v_%v_%v",conf.Addr,conf.DB,conf.Username,conf.Network)
}
