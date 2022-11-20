package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"strconv"
	"time"
)

func GetBaseConfig(c framework.Container) *contract.RedisConfig {
	logService := c.MustMake(contract.LogKey).(contract.Log)
	config := &contract.RedisConfig{Options: &redis.Options{}}
	opt := WithConfigPath("redis")
	err := opt(c,config)
	if err != nil {
		logService.Error(context.Background(),"parse cache conf err",nil)
		return nil
	}
	return config
}

func WithConfigPath(configPath string) contract.RedisOption {
	return func(c framework.Container, config *contract.RedisConfig) error {
		configService := c.MustMake(contract.ConfigKey).(contract.Config)
		conf := configService.GetStringMapString(configPath)
		if host, ok := conf["host"]; ok {
			if port, ok1 := conf["port"]; ok1 {
				config.Addr = host + ":" + port
			}
			fmt.Println(config.Addr)
		}
		if db, ok := conf["db"]; ok {
			t, err := strconv.Atoi(db)
			if err != nil {
				return err
			}
			config.DB = t
		}
		if username,ok := conf["username"]; ok {
			config.Username = username
		}
		if password, ok :=  conf["password"]; ok {
			config.Password = password
		}
		if timeout, ok := conf["timeout"]; ok {
			t, err := time.ParseDuration(timeout);
			if err != nil {
				return err
			}
			config.DialTimeout = t
		}

		if timeout, ok := conf["read_timeout"]; ok {
			t,err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.ReadTimeout = t
		}

		if timeout, ok := conf["write_timeout"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.WriteTimeout = t
		}

		if cnt, ok := conf["conn_min_idle"]; ok {
			t, err := strconv.Atoi(cnt)
			if err != nil {
				return err
			}
			config.MinIdleConns = t
		}

		if max, ok := conf["conn_max_open"]; ok {
			t, err := strconv.Atoi(max)
			if err != nil {
				return err
			}
			config.PoolSize = t
		}

		if timeout, ok := conf["conn_max_lifetime"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.MaxConnAge = t
		}

		if timeout, ok := conf["conn_max_idletime"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.IdleTimeout = t
		}
		return nil
	}
}

///表示自行配置redis的配置信息
func WithRedisConfig(f func(options *contract.RedisConfig)) contract.RedisOption {
	return func(container framework.Container, config *contract.RedisConfig) error {
		f(config)
		return nil
	}
}
