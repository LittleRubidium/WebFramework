package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/erikdubbelboer/gspt"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/cobra"
	"github.com/gohade/hade/framework/contract"
	"github.com/gohade/hade/framework/util"
	"github.com/sevlyar/go-daemon"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

var appAddress = ""
var appDaemon = false

func initAppCommand() *cobra.Command {
	appStartCommand.Flags().BoolVarP(&appDaemon,"daemon","d",false,"start app daemon")
	appStartCommand.Flags().StringVar(&appAddress,"address","","设置app启动端口，默认为8888")

	appCommand.AddCommand(appStartCommand)
	appCommand.AddCommand(appRestartCommand)
	appCommand.AddCommand(appStopCommand)
	appCommand.AddCommand(appStateCommand)
	return appCommand
}

var appCommand = &cobra.Command{
	Use: "app",
	Short: "业务应用控制命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

//启动AppServer，这个函数会将当前goroutine阻塞
func startAppServe(server *http.Server, c framework.Container) error {
	go server.ListenAndServe()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	closeWait := 5
	configServer := c.MustMake(contract.ConfigKey).(contract.Config)
	if configServer.IsExist("app.close_wait") {
		closeWait = configServer.GetInt("app.close_wait")
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(closeWait) *time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	return nil

}

var appStartCommand = &cobra.Command{
	Use: "start",
	Short: "启动一个服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		kernelService := container.MustMake(contract.KernelKey).(contract.Kernel)
		core := kernelService.HttpEngine()

		if appAddress == "" {
			envServer := container.MustMake(contract.EnvKey).(contract.Env)
			if envServer.Get("ADDRESS") != "" {
				appAddress = envServer.Get("ADDRESS")
			}else {
				configService := container.MustMake(contract.ConfigKey).(contract.Config)
				if configService.IsExist("app.address") {
					appAddress = configService.GetString("app.address")
				}else {
					appAddress = ":8888"
				}
			}
		}
		//创建一个Server服务
		server := &http.Server{
			Handler: core,
			Addr: appAddress,
		}

		appService := container.MustMake(contract.AppKey).(contract.App)

		pidFolder := appService.RuntimeFolder()
		if !util.Exists(pidFolder) {
			if err := os.MkdirAll(pidFolder,os.ModePerm); err != nil {
				return err
			}
		}
		serverPidFile := filepath.Join(pidFolder,"app.pid")
		logFolder := appService.LogFolder()
		if !util.Exists(logFolder) {
			if err := os.MkdirAll(logFolder,os.ModePerm); err != nil {
				return err
			}
		}
		serverLogFile := filepath.Join(logFolder,"app.log")
		currentFolder := util.GetExecDirectory()
		//daemon模式
		if appDaemon {
			//创建一个context
			cntxt := &daemon.Context{
				//设置pid文件
				PidFileName: serverPidFile,
				PidFilePerm: 0664,
				LogFileName: serverLogFile,
				LogFilePerm: 0640,
				WorkDir:     currentFolder,
				Umask:       027,
				Args:        []string{"", "app", "start", "--daemon=true"},
			}
			d, err := cntxt.Reborn()
			if err != nil {
				return err
			}
			if d != nil {
				fmt.Println("app启动成功，pid:",d.Pid)
				fmt.Println("日志文件:",serverLogFile)
				return nil
			}
			defer cntxt.Release()
			fmt.Println("daemon started")
			gspt.SetProcTitle("hade app")
			if err := startAppServe(server,container); err != nil {
				fmt.Println(err)
			}
			return nil
		}

		//非daemon模式，直接执行
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID]",content)
		err := ioutil.WriteFile(serverPidFile,[]byte(content),0664)
		if err != nil {
			return err
		}
		gspt.SetProcTitle("hade app")
		fmt.Println("app serve url:",appAddress)
		if err := startAppServe(server, container); err != nil {
			fmt.Println(err)
		}
		return nil
	},
}

//重启一个服务
var appRestartCommand = &cobra.Command{
	Use: "restart",
	Short: "重新启动一个服务",
	RunE: func(c *cobra.Command, args []string) error{
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		serverPidFile := filepath.Join(appService.RuntimeFolder(),"app.pid")

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if content != nil && len(content) != 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				if err := syscall.Kill(pid,syscall.SIGTERM); err != nil {
					return err
				}

				//获取closeWait
				closeWait := 5
				configService := container.MustMake(contract.ConfigKey).(contract.Config)
				if configService.IsExist("app.close_wait") {
					closeWait = configService.GetInt("app.close_wait")
				}

				for i := 0;i < closeWait * 2;i++ {
					if !util.CheckProcessExist(pid) {
						break
					}
					time.Sleep(1 * time.Second)
				}

				if util.CheckProcessExist(pid) {
					fmt.Println("结束进程失败:" + strconv.Itoa(pid),"请查看原因")
					return errors.New("结束进程失败")
				}
				if err := ioutil.WriteFile(serverPidFile,[]byte{},0664); err != nil {
					return err
				}
				fmt.Println("结束进程成功:" + strconv.Itoa(pid))
			}
		}
		appDaemon = true
		return appStartCommand.RunE(c,args)
	},
}

// 停止一个已经启动的app服务
var appStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "停止一个已经启动的app服务",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// GetPid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if content != nil && len(content) != 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			// 发送SIGTERM命令
			if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
				return err
			}
			if err := ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
				return err
			}
			fmt.Println("停止进程:", pid)
		}
		return nil
	},
}

// 获取启动的app的pid
var appStateCommand = &cobra.Command{
	Use:   "state",
	Short: "获取启动的app的pid",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// 获取pid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				fmt.Println("app服务已经启动, pid:", pid)
				return nil
			}
		}
		fmt.Println("没有app服务存在")
		return nil
	},
}