package cobra

import (
	"github.com/gohade/hade/framework/contract"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

//AddDistributedCronCommand实现一个分布式定时器 
//serviceName这个服务的唯一名字，不允许带有空格
//spec具体的执行时间
//cmd具体的执行命令
//holdTime表示如果被选择的持续有效时间，也就是持有锁的时间
func (c *Command) AddDistributedCronCommand(serviceName string, spec string, cmd *Command,holdTime time.Duration)  {
	root := c.Root()
	if root.Cron == nil {
		root.Cron = cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))
		root.CronSpecs = []CronSpec{}
	}
	root.CronSpecs = append(root.CronSpecs,CronSpec{
		Type: "distributed-cron",
		Cmd: cmd,
		Spec: spec,
		ServiceName: serviceName,
	})
	appService := root.GetContainer().MustMake(contract.AppKey).(contract.App)
	distributeService := root.GetContainer().MustMake(contract.DistributedKey).(contract.Distributed)
	appID := appService.AppID()
	var cronCmd Command
	ctx := root.Context()
	cronCmd = *cmd
	cronCmd.args = []string{}
	cronCmd.SetParentNull()
	root.Cron.AddFunc(spec, func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		selectedAppId,err := distributeService.Select(serviceName,appID,holdTime)
		if err != nil {
			return
		}
		if selectedAppId != appID {
			return
		}
		err = cronCmd.ExecuteContext(ctx)
		if err != nil {
			log.Println(err)
		}
	})
}