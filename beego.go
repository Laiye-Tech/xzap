package xzap

import "context"

type BeegoLog struct {
}

func (b BeegoLog) Error(format string, v ...interface{}) {
	beeLogger.Sugar().Errorf(format, v...)
}

func (b BeegoLog) Warning(format string, v ...interface{}) {
	beeLogger.Sugar().Warnf(format, v...)
}

func (b BeegoLog) Informational(format string, v ...interface{}) {
	beeLogger.Sugar().Infof(format, v...)
}

func (b BeegoLog) Debug(format string, v ...interface{}) {
	beeLogger.Sugar().Debugf(format, v...)

}
func (b BeegoLog) Warn(format string, v ...interface{}) {
	beeLogger.Sugar().Warnf(format, v...)
}

func (b BeegoLog) Info(format string, v ...interface{}) {
	beeLogger.Sugar().Infof(format, v...)
}

func (b BeegoLog) ErrorCtx(ctx context.Context, format string, v ...interface{}) {
	with(ctx, beeLogger).Sugar().Errorf(format, v...)
}

func (b BeegoLog) WarningCtx(ctx context.Context, format string, v ...interface{}) {
	with(ctx, beeLogger).Sugar().Warnf(format, v...)
}

func (b BeegoLog) InformationalCtx(ctx context.Context, format string, v ...interface{}) {
	with(ctx, beeLogger).Sugar().Infof(format, v...)
}

func (b BeegoLog) DebugCtx(ctx context.Context, format string, v ...interface{}) {
	with(ctx, beeLogger).Sugar().Debugf(format, v...)
}

func (b BeegoLog) WarnCtx(ctx context.Context, format string, v ...interface{}) {
	with(ctx, beeLogger).Sugar().Warnf(format, v...)
}

func (b BeegoLog) InfoCtx(ctx context.Context, format string, v ...interface{}) {
	with(ctx, beeLogger).Sugar().Infof(format, v...)
}

type BeegoCtxLog struct {
	ctx context.Context
}

func (b BeegoCtxLog) Error(format string, v ...interface{}) {
	with(b.ctx, beeLogger).Sugar().Errorf(format, v...)
}

func (b BeegoCtxLog) Warning(format string, v ...interface{}) {
	with(b.ctx, beeLogger).Sugar().Warnf(format, v...)
}

func (b BeegoCtxLog) Informational(format string, v ...interface{}) {
	with(b.ctx, beeLogger).Sugar().Infof(format, v...)
}

func (b BeegoCtxLog) Debug(format string, v ...interface{}) {
	with(b.ctx, beeLogger).Sugar().Debugf(format, v...)

}
func (b BeegoCtxLog) Warn(format string, v ...interface{}) {
	with(b.ctx, beeLogger).Sugar().Warnf(format, v...)
}

func (b BeegoCtxLog) Info(format string, v ...interface{}) {
	with(b.ctx, beeLogger).Sugar().Infof(format, v...)
}
