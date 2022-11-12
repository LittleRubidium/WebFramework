package contract

const AppKey = "hade:app"

type App interface {
	//Version 定义当前版本
	Version() string
	//BaseFolder定义基础地址
	BaseFolder() string
	//ConfigFolder定义 配置文件的路径
	ConfigFolder() string
	//LogFolder定义日志所在的路径
	LogFolder() string
	//ProviderFolder定义业务自己的提供者地址
	ProviderFolder() string
	//MiddlewareFolder定义业务自己的中间件
	MiddlewareFolder() string
	//CommandFolder定义业务自己的命令
	CommandFolder() string
	//RuntimeFolder定义业务运行中间态信息
	RuntimeFolder() string
	//TestFolder存放测试所需要的信息
	TestFolder() string
}