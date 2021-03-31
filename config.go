package xzap

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AsyncParams struct {
	//default true
	Async bool
	//default 500ms
	FlushInternal time.Duration
	//default 100
	FlushLine int
}

type config struct {
	//beego like format config
	zapConfig *zap.Config
	//default true
	daily bool
	//max size in megabytes
	//default 256M
	maxSize int
	//default 7
	maxAge int
	//default info
	level level
	//atomic logging level, can be modified runtime.
	atomLevel zap.AtomicLevel
	//
	async AsyncParams
}

func defaultConfig() *config {
	return &config{
		zapConfig: zapConfig(),
		daily:     true,
		maxSize:   256,
		maxAge:    7,
		level:     Info,
		atomLevel: atomLevel,
		async: AsyncParams{
			Async:         true,
			FlushInternal: time.Millisecond * 500,
			FlushLine:     100,
		},
	}
}

type level string

const Debug level = "Debug"
const Info level = "Info"
const Warn level = "Warn"
const Error level = "Error"

func (l level) ZapLevel() zapcore.Level {
	switch l {
	case "":
		return zapcore.InfoLevel
	case Debug:
		return zapcore.DebugLevel
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Error:
		return zapcore.ErrorLevel
	default:
		fmt.Println("un know level, use default level info")
		return zapcore.InfoLevel
	}
}
