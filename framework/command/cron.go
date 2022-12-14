package command

import (
	"fmt"
	"github.com/erikdubbelboer/gspt"
	"github.com/gohade/hade/framework/cobra"
	"github.com/gohade/hade/framework/contract"
	"github.com/gohade/hade/framework/util"
	"github.com/sevlyar/go-daemon"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

var cronDaemon = false

func initCronCommand() *cobra.Command {
	cronStartCommand.Flags().BoolVarP(&cronDaemon,"daemon","d",false,"start serve cron")
	cronCommand.AddCommand(cronStartCommand)
	cronCommand.AddCommand(cronStateCommand)
	cronCommand.AddCommand(cronStopCommand)
	cronCommand.AddCommand(cronListCommand)
	cronCommand.AddCommand(cronRestartCommand)
	return cronCommand
}

var cronCommand = &cobra.Command{
	Use: "cron",
	Short: "定时任务相关命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

var cronListCommand = &cobra.Command{
	Use: "list",
	Short: "列出所有的定时任务",
	RunE: func(cmd *cobra.Command, args []string) error{
		cronSpecs := cmd.Root().CronSpecs
		ps := [][]string{}
		for _, cronSpec := range cronSpecs {
			line := []string{cronSpec.Type,cronSpec.Spec,cronSpec.Cmd.Use, cronSpec.Cmd.Short, cronSpec.ServiceName}
			ps = append(ps,line)
		}
		util.PrettyPrint(ps)
		return nil
	},
}

var cronStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动cron常驻进程",
	RunE: func(c *cobra.Command, args []string) error {
		//获取容器
		container := c.GetContainer()
		//获取容器中的app服务
		appService := container.MustMake(contract.AppKey).(contract.App)

		//设置cron的日志地址和进程id地址
		pidFolder := appService.RuntimeFolder()
		serverPidFolder := filepath.Join(pidFolder,"cron.pid")
		logFolder := appService.LogFolder()
		serverLogFile := filepath.Join(logFolder,"cron.log")
		currentFolder := appService.BaseFolder()

		if cronDaemon {
			//创建一个context
			cntxt := daemon.Context{
				PidFileName: serverPidFolder,
				PidFilePerm: 0664,
				LogFileName: serverLogFile,
				LogFilePerm: 0640,
				WorkDir: currentFolder,
				Umask: 027,
				Args: []string{"","cron","start","--daemon=true"},
			}
			//启动子进程，d不为空表示当前是父进程，d为空表示当前是子进程
			d, err := cntxt.Reborn()
			if err != nil {
				return err
			}
			if d != nil {
				fmt.Println("cron serve started, pid: ",d.Pid)
				fmt.Println("log file: ",serverLogFile)
				return nil
			}

			//子进程Cron.Run
			defer cntxt.Release()
			fmt.Println("daemon started")
			gspt.SetProcTitle("hade cron")
			c.Root().Cron.Run()
			return nil
		}
		fmt.Println("start cron job")
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID]",content)
		err := ioutil.WriteFile(serverPidFolder,[]byte(content),0664)
		if err != nil {
			return err
		}

		gspt.SetProcTitle("hade cron")
		c.Root().Cron.Run()
		return nil
	},
}

var cronRestartCommand = &cobra.Command{
	Use: "restart",
	Short: "重启cron常驻进程",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		serverPidFile := filepath.Join(appService.RuntimeFolder(),"cron pid")
		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}
		if content != nil && len(content) > 0 {
			pid,err:= strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				if err := syscall.Kill(pid,syscall.SIGTERM); err != nil {
					return err
				}

				//check process closed
				for i := 0;i < 10; i++ {
					if !util.CheckProcessExist(pid) {
						break
					}
					time.Sleep(1 * time.Second)
				}
				fmt.Println("kill process: " + strconv.Itoa(pid))
			}
		}
		cronDaemon = true
		return cronStartCommand.RunE(c,args)
	},
}

var cronStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "停止cron常驻进程",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// GetPid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
				return err
			}
			if err := ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
				return err
			}
			fmt.Println("stop pid:", pid)
		}
		return nil
	},
}

var cronStateCommand = &cobra.Command{
	Use:   "state",
	Short: "cron常驻进程状态",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// GetPid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")

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
				fmt.Println("cron server started, pid:", pid)
				return nil
			}
		}
		fmt.Println("no cron server start")
		return nil
	},
}