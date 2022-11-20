package contract

import "time"

const (
	ConfigKey = "hade:config"
)

type Config interface {
	//检查一个属性是否存在
	IsExist(key string) bool
	//获取一个属性
	Get(key string) interface{}
	GetBool(key string) bool
	GetInt(key string) int
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetString(key string) string
	GetIntSlice(key string) []int
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	Load(key string, val interface{}) error
}
