package contract

import (
	"context"
	"io"
	"time"
)

const (
	LogKey = "hade:log"
)

type LogLevel int32

//日志级别
const (
	UnKnowLevel LogLevel = iota
	PanicLevel
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

//CtxFielder 定义了从context中获取信息的方法
type CtxFielder func(ctx context.Context) map[string]interface{}

//Formatter定义了将日志信息组织成字符串的通用方法
type Formatter func(level LogLevel, t time.Time, msg string, field map[string]interface{}) ([]byte, error)

//Log define interface for log
type Log interface {
	//从context中获取上下文字段field
	SetCtxFielder(handler CtxFielder)
	//设置输出格式
	SetFormatter(handler Formatter)
	//设置输出管道
	SetOutput(out io.Writer)

	Panic(ctx context.Context, msg string, field map[string]interface{})
	Fatal(ctx context.Context, msg string, fields map[string]interface{})
	Error(ctx context.Context, msg string, fields map[string]interface{})
	Warn(ctx context.Context, msg string, fields map[string]interface{})
	Info(ctx context.Context, msg string, fields map[string]interface{})
	Debug(ctx context.Context, msg string, fields map[string]interface{})
	Trace(ctx context.Context, msg string, fields map[string]interface{})
	//设置日志级别
	SetLevel(level LogLevel)
}
