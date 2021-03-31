package xzap

import (
	"context"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"git.laiye.com/laiye-backend-repos/go-utils/xzap/lumberjack"
)

var logger *zap.Logger
var beeLogger *zap.Logger
var atomLevel = zap.NewAtomicLevel() // 管理运行时日志级别

func init() {
	logger, _ = zap.NewDevelopmentConfig().Build()
	beeLogger = logger.WithOptions(zap.AddCallerSkip(1))
	zap.ReplaceGlobals(logger)
}

var (
	enc zapcore.Encoder
	ws  zapcore.WriteSyncer
)

var inited int32
var once sync.Once

func InitLog(fileName string, opts ...Option) {
	if atomic.LoadInt32(&inited) > 0 {
		panic("xzap already inited")
	}
	atomic.AddInt32(&inited, 1)
	once.Do(func() {
		go startRuntimeLevelServer()
		initLog(fileName, opts...)
	})
}

func New(fileName string, opts ...Option) *zap.Logger {
	c := defaultConfig()
	for _, f := range opts {
		f.Apply(c)
	}
	ws = lumberjack.NewLogger(&lumberjack.Logger{
		AsyncParams: c.async,
		Filename:    fileName,
		MaxSize:     c.maxSize,
		MaxAge:      c.maxAge, // days
		Daily:       c.daily,
		LocalTime:   true,
	})
	cfg := c.zapConfig
	enc = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	c.atomLevel.SetLevel(c.level.ZapLevel())
	core := zapcore.NewCore(
		enc,
		ws,
		c.atomLevel,
	)

	logger = zap.New(core, zap.AddCaller())
	return logger
}

func initLog(fileName string, opts ...Option) {
	logger = New(fileName, opts...)
	beeLogger = logger.WithOptions(zap.AddCallerSkip(1))
	zap.ReplaceGlobals(logger)
}

func With(ctx context.Context, lg *zap.Logger) *zap.Logger {
	return with(ctx, lg)
}

func with(ctx context.Context, lg *zap.Logger) *zap.Logger {
	if ctx == nil {
		return lg
	}
	newLogger := lg

	if sc := getTraceSpanContext(ctx); sc != nil {
		newLogger = newLogger.With(
			zap.String("traceID", sc.TraceID().String()),
			zap.Bool("isSampled", sc.IsSampled()),
		)
	}
	if debug(ctx) {
		if enc != nil && ws != nil {
			newLogger = newLogger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
				return zapcore.NewCore(enc, ws, zapcore.DebugLevel)
			}))
		}
	}
	if v := msgID(ctx); v != 0 {

		newLogger = newLogger.With(zap.Uint64("msgID", v))
	} else {
		msgID := ctx.Value("msgID")
		if msgID != nil {
			switch v := msgID.(type) {
			case string:
				newLogger = newLogger.With(zap.String("msgID", v))
			case uint64:
				if v != 0 {
					newLogger = newLogger.With(zap.Uint64("msgID", v))
				}
			case int64:
				if v != 0 {
					newLogger = newLogger.With(zap.Int64("msgID", v))
				}
			default:
				newLogger = newLogger.With(zap.Any("msgID", v))
			}
		}

	}
	return newLogger
}

func Logger(ctx context.Context) *zap.Logger {
	return with(ctx, logger)
}

func Sugar(ctx context.Context) *zap.SugaredLogger {
	return Logger(ctx).Sugar()
}

func Beego() BeegoLog {
	return BeegoLog{}
}

func BeegoCtx(ctx context.Context) BeegoCtxLog {
	return BeegoCtxLog{ctx: ctx}
}
