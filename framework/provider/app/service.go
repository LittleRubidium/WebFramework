package app

import (
	"errors"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/util"
	"github.com/google/uuid"
	"path/filepath"
)

//HadeApp
type HadeApp struct {
	container framework.Container  //服务容器
	baseFolder string //基础路径
	appID string
}

//实现版本
func (h HadeApp) Version() string {
	return "0.0.1"
}

//基础目录，可根据开发或运行时调整
func (h HadeApp) BaseFolder() string {
	if h.baseFolder != "" {
		return h.baseFolder
	}

	//没有指定路径，则使用默认路径
	return util.GetExecDirectory()
}

// ConfigFolder  表示配置文件地址
func (h HadeApp) ConfigFolder() string {
	return filepath.Join(h.BaseFolder(), "config")
}

// LogFolder 表示日志存放地址
func (h HadeApp) LogFolder() string {
	return filepath.Join(h.StorageFolder(), "log")
}

func (h HadeApp) HttpFolder() string {
	return filepath.Join(h.BaseFolder(), "http")
}

func (h HadeApp) ConsoleFolder() string {
	return filepath.Join(h.BaseFolder(), "console")
}

func (h HadeApp) StorageFolder() string {
	return filepath.Join(h.BaseFolder(), "storage")
}

// ProviderFolder 定义业务自己的服务提供者地址
func (h HadeApp) ProviderFolder() string {
	return filepath.Join(h.BaseFolder(), "provider")
}

// MiddlewareFolder 定义业务自己定义的中间件
func (h HadeApp) MiddlewareFolder() string {
	return filepath.Join(h.HttpFolder(), "middleware")
}

// CommandFolder 定义业务定义的命令
func (h HadeApp) CommandFolder() string {
	return filepath.Join(h.ConsoleFolder(), "command")
}

// RuntimeFolder 定义业务的运行中间态信息
func (h HadeApp) RuntimeFolder() string {
	return filepath.Join(h.StorageFolder(), "runtime")
}

// TestFolder 定义测试需要的信息
func (h HadeApp) TestFolder() string {
	return filepath.Join(h.BaseFolder(), "test")
}

func (h HadeApp) AppID() string {
	return h.appID
}

func NewHadeApp(params []interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil,errors.New("param error")
	}
	appId := uuid.New().String()
	//有两个参数，一个是容器，一个是baseFolder
	container := params[0].(framework.Container)
	baseFolder := params[1].(string)
	return &HadeApp{baseFolder: baseFolder,container: container,appID: appId},nil
}

