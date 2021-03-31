package xzap

import (
	"fmt"
	"reflect"
	"sync"
	"time"
	"unsafe"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"git.laiye.com/laiye-backend-repos/go-utils/xzap/tfmt"
)

func zapConfig() *zap.Config {
	c := zap.NewProductionConfig()
	c.Encoding = "console"
	c.EncoderConfig.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		f := func(level zapcore.Level) string {
			switch level {
			case zapcore.DebugLevel:
				return "[D]"
			case zapcore.InfoLevel:
				return "[I]"
			case zapcore.WarnLevel:
				return "[W]"
			case zapcore.ErrorLevel:
				return "[E]"
			case zapcore.DPanicLevel:
				return "[C]"
			case zapcore.PanicLevel:
				return "[A]"
			case zapcore.FatalLevel:
				return "[M]"
			default:
				return fmt.Sprintf("level(%d)", level)
			}
		}
		encoder.AppendString(f(level))
	}

	c.EncoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		t, _, _ := tfmt.FormatTimeHeader(time, nil)
		s := *(*string)(unsafe.Pointer(&t))
		encoder.AppendString(s)
	}

	c.EncoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
		ph, ok := getCaller(caller.File, caller.Line)
		if !ok {
			ph = "[" + caller.TrimmedPath() + "]"
			setCaller(caller.File, caller.Line, ph)
		}
		encoder.AppendString(ph)
	}
	return &c
}

var stringPool sync.Map

func hash(file string, line int) uint64 {
	return MemHashString(file)<<10 + uint64(line)
}

func getCaller(file string, line int) (string, bool) {
	s, ok := stringPool.Load(hash(file, line))
	if !ok {
		return "", false
	}
	return s.(string), ok
}

func setCaller(file string, line int, caller string) {
	stringPool.LoadOrStore(hash(file, line), caller)
}

func MemHashString(str string) uint64 {
	ss := (*reflect.StringHeader)(unsafe.Pointer(&str))
	return uint64(memhash(unsafe.Pointer(ss.Data), 0, uintptr(ss.Len)))
}

//go:noescape
//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, h, s uintptr) uintptr
