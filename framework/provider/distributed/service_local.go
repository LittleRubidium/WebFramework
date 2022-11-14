package distributed

import (
	"errors"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type LocalDistributedService struct {
	container framework.Container
}

func NewLocalDistributedService(params []interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("param error")
	}

	//有两个参数，一个是容器，一个是baseFolder
	container := params[0].(framework.Container)
	return &LocalDistributedService{container: container},nil
}

//分布式选择器
func (l LocalDistributedService) Select(serviceName string, appID string, holdtime time.Duration) (selectAppID string, err error) {
	appService := l.container.MustMake(contract.AppKey).(contract.App)
	runtimeFolder := appService.RuntimeFolder()
	lockFile := filepath.Join(runtimeFolder,"distribute_" + serviceName)

	//打开文件锁
	lock, err := os.OpenFile(lockFile, os.O_RDWR | os.O_CREATE,0666)
	if err != nil {
		return "", err
	}
	//尝试独占文件锁
	err = syscall.Flock(int(lock.Fd()),syscall.LOCK_EX | syscall.LOCK_NB)
	//抢不到文件锁
	if err != nil {
		selectAppIDByt, err := ioutil.ReadAll(lock)
		if err != nil {
			return "", err
		}
		return string(selectAppIDByt),err
	}

	//在一段时间内，选举有效，其他节点在这段时间不能再进行抢占
	go func() {
		defer func() {
			syscall.Flock(int(lock.Fd()),syscall.LOCK_UN)
			lock.Close()
			os.Remove(lockFile)
		}()

		timer := time.NewTimer(holdtime)
		<- timer.C
	}()
	if _, err := lock.WriteString(appID); err != nil {
		return "", err
	}
	return appID, nil
}