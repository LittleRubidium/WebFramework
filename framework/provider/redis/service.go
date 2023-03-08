package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"sync"
)

type HadeRedis struct {
	container framework.Container
	clients   map[string]*redis.Client

	lock *sync.RWMutex
}

func NewHadeRedis(params []interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	clients := make(map[string]*redis.Client)
	lock := &sync.RWMutex{}
	return &HadeRedis{
		container: container,
		clients:   clients,
		lock:      lock,
	}, nil
}

func (r *HadeRedis) GetClient(option ...contract.RedisOption) (*redis.Client, error) {
	//读取默认配置
	config := GetBaseConfig(r.container)
	//option对opt进行修改
	for _, opt := range option {
		if err := opt(r.container, config); err != nil {
			return nil, err
		}
	}

	//如果最终的config没有设置dsn，就设置dsn
	key := config.UniqKey()
	//fmt.Println(key)

	//判断是否已经实例化redis.Client，就生dsn
	r.lock.RLock()
	if db, ok := r.clients[key]; ok {
		r.lock.RUnlock()
		return db, nil
	}
	r.lock.RUnlock()
	r.lock.Lock()
	defer r.lock.Unlock()

	//实例化gorm.DB
	client := redis.NewClient(config.Options)

	r.clients[key] = client
	return client, nil
}
