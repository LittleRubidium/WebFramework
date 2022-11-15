package log

import (
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"github.com/gohade/hade/framework/provider/log/formatter"
	"github.com/gohade/hade/framework/provider/log/services"
	"io"
	"strings"
)

type HadeLogServiceProvider struct {
	Driver string
	Level contract.LogLevel
	Formatter contract.Formatter
	CtxFielder contract.CtxFielder
	Output io.Writer
}

func (log *HadeLogServiceProvider) Register(c framework.Container) framework.NewInstance {
	if log.Driver == "" {
		tcs,err := c.Make(contract.ConfigKey)
		if err != nil {
			return services.NewHadeConsoleLog
		}

		cs := tcs.(contract.Config)
		log.Driver = strings.ToLower(cs.GetString("log.Driver"))
	}

	//根据 driver的配置项确定
	switch log.Driver {
	case "single":
		return services.NewHadeSingleLog
	case "rotate":
		return services.NewHadeRotateLog
	case "console":
		return services.NewHadeConsoleLog
	case "custom":
		return services.NewHadeCustomLog
	default:
		return services.NewHadeConsoleLog
	}
}

func (log *HadeLogServiceProvider) Boot(c framework.Container) error {
	return nil
}

func (log *HadeLogServiceProvider) IsDefer() bool {
	return false
}

func (log *HadeLogServiceProvider) Params(c framework.Container) []interface{} {
	configService := c.MustMake(contract.ConfigKey).(contract.Config)

	//设置参数formatter
	if log.Formatter == nil {
		log.Formatter = formatter.TextFormatter
		if configService.IsExist("log.formatter") {
			v := configService.GetString("log.formatter")
			if v == "json" {
				log.Formatter = formatter.JsonFormatter
			}else if v == "text" {
				log.Formatter = formatter.TextFormatter
			}
		}
	}

	if log.Level == contract.UnKnowLevel {
		log.Level = contract.InfoLevel
		if configService.IsExist("log.level") {
			log.Level = logLevel(configService.GetString("log.level"))
		}
	}
	return []interface{}{c, log.Level, log.CtxFielder, log.Formatter, log.Output}
}

func (log *HadeLogServiceProvider) Name() string {
	return contract.LogKey
}

func logLevel(config string) contract.LogLevel {
	switch strings.ToLower(config) {
	case "panic":
		return contract.PanicLevel
	case "fatal":
		return contract.FatalLevel
	case "error":
		return contract.ErrorLevel
	case "warn":
		return contract.WarnLevel
	case "info":
		return contract.InfoLevel
	case "trace":
		return contract.TraceLevel
	}
	return contract.UnKnowLevel
}