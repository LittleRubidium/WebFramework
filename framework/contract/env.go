package contract

const (
	//生产环境
	EnvProduction = "production"
	//测试环境
	EnvTesting = "testing"
	//开发环境
	EnvDevelopment = "development"
	//环境变量服务字符串凭证
	EnvKey = "hade:env"
)

//环境变量服务
type Env interface {
	//AppEnv获取当前的环境
	AppEnv() string
	//判断一个环境变量是否有被设置
	IsExist(string) bool
	//Get获取某个环境变量，如果没有设置则，返回""
	Get(string) string
	//获取所有环境变量，.env和运行环境变量融合后的结果
	All() map[string]string
}
