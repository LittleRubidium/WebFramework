package app

import (
	"errors"
	"flag"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/util"
	"path/filepath"
)

//HadeApp
type HadeApp struct {
	container framework.Container  //服务容器
	baseFolder string //基础路径
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

	//如果没有baseFolder,则设置
	var baseFolder string
	flag.StringVar(&baseFolder,"base_folder","","base_folder参数，默认为当前路径")
	flag.Parse()
	if baseFolder != "" {
		return baseFolder
	}
	//没有指定路径，则使用默认路径
	return util.GetExecDirectory()
}

func (h HadeApp) ConfigFolder() string {
	return filepath.Join(h.BaseFolder(),"config")
}

func (h HadeApp) LogFolder() string {
	return filepath.Join(h.BaseFolder(),"log")
}

func (h HadeApp) HttpFolder() string {
	return filepath.Join(h.BaseFolder(),"http")
}

func (h HadeApp) StorageFolder() string {
	return filepath.Join(h.BaseFolder(),"storage")
}

func (h HadeApp) ProviderFolder() string {
	return filepath.Join(h.BaseFolder(),"provider")
}

func (h HadeApp) MiddlewareFolder() string {
	return filepath.Join(h.BaseFolder(),"middleware")
}

func (h HadeApp) CommandFolder() string {
	return filepath.Join(h.BaseFolder(),"command")
}

func (h HadeApp) RuntimeFolder() string {
	return filepath.Join(h.BaseFolder(),"runtime")
}

func (h HadeApp) TestFolder() string {
	return filepath.Join(h.BaseFolder(),"test")
}

func NewHadeApp(params []interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil,errors.New("param error")
	}

	//有两个参数，一个是容器，一个是baseFolder
	container := params[0].(framework.Container)
	baseFolder := params[1].(string)
	return &HadeApp{baseFolder: baseFolder,container: container},nil
}

