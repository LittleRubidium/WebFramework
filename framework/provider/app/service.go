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
	configMap map[string]string
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
func (a HadeApp) ConfigFolder() string {
	if val, ok := a.configMap["config_folder"]; ok {
		return val
	}
	return filepath.Join(a.BaseFolder(), "config")
}

// LogFolder 表示日志存放地址
func (a HadeApp) LogFolder() string {
	if val, ok := a.configMap["log_folder"]; ok {
		return val
	}
	return filepath.Join(a.StorageFolder(), "log")
}

func (a HadeApp) HttpFolder() string {
	if val, ok := a.configMap["http_folder"]; ok {
		return val
	}
	return filepath.Join(a.BaseFolder(), "http")
}

func (a HadeApp) ConsoleFolder() string {
	if val, ok := a.configMap["console_folder"]; ok {
		return val
	}
	return filepath.Join(a.BaseFolder(), "console")
}

func (a HadeApp) StorageFolder() string {
	if val, ok := a.configMap["storage_folder"]; ok {
		return val
	}
	return filepath.Join(a.BaseFolder(), "storage")
}

// ProviderFolder 定义业务自己的服务提供者地址
func (a HadeApp) ProviderFolder() string {
	if val, ok := a.configMap["provider_folder"]; ok {
		return val
	}
	return filepath.Join(a.BaseFolder(), "provider")
}

// MiddlewareFolder 定义业务自己定义的中间件
func (a HadeApp) MiddlewareFolder() string {
	if val, ok := a.configMap["middleware_folder"]; ok {
		return val
	}
	return filepath.Join(a.HttpFolder(), "middleware")
}

// CommandFolder 定义业务定义的命令
func (a HadeApp) CommandFolder() string {
	if val, ok := a.configMap["command_folder"]; ok {
		return val
	}
	return filepath.Join(a.ConsoleFolder(), "command")
}

// RuntimeFolder 定义业务的运行中间态信息
func (a HadeApp) RuntimeFolder() string {
	if val, ok := a.configMap["runtime_folder"]; ok {
		return val
	}
	return filepath.Join(a.StorageFolder(), "runtime")
}

// TestFolder 定义测试需要的信息
func (a HadeApp) TestFolder() string {
	if val, ok := a.configMap["test_folder"]; ok {
		return val
	}
	return filepath.Join(a.BaseFolder(), "test")
}

func (a HadeApp) AppID() string {
	return a.appID
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

func (a *HadeApp) LoadAppConfig(kv map[string]string) {
	for key, val := range kv {
		a.configMap[key] = val
	}
}

