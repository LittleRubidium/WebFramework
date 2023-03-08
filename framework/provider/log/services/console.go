package services

import (
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"os"
)

type HadeConsoleLog struct {
	HadeLog
}

func NewHadeConsoleLog(params []interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	log := &HadeConsoleLog{}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
	log.c = c
	return log, nil
}
