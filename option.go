package xzap

import "go.uber.org/zap"

type Option interface {
	Apply(*config)
}

type funcOption func(c *config)

func (f funcOption) Apply(c *config) {
	f(c)
}

func WithDaily(daily bool) Option {
	return withDaily(daily)
}

type withDaily bool

func (w withDaily) Apply(c *config) {
	c.daily = bool(w)
}

func WithMaxSize(maxSize int) Option {
	return withMaxSize(maxSize)
}

type withMaxSize int

func (w withMaxSize) Apply(c *config) {
	if w == 0 {
		return
	}
	c.maxSize = int(w)
}

func WithMaxAge(maxAge int) Option {
	return withMaxAge(maxAge)
}

type withMaxAge int

func (w withMaxAge) Apply(c *config) {
	if w == 0 {
		return
	}
	c.maxAge = int(w)
}

func WithAsyncParams(s AsyncParams) Option {
	return funcOption(func(c *config) {
		c.async.Async = s.Async
		if s.FlushInternal != 0 {
			c.async.FlushInternal = s.FlushInternal
		}
		if s.FlushLine != 0 {
			c.async.FlushLine = s.FlushLine
		}
	})
}

func WithLevel(level level) Option {
	return withLevel(level)
}

type withLevel level

func (w withLevel) Apply(c *config) {
	c.level = level(w)
}

//WithAtomLevel 如果使用此选项，默认的运行时修改日志等级会失效，使用者可自行管理运行时日志等级
func WithAtomLevel(level zap.AtomicLevel) Option {
	return funcOption(func(c *config) {
		c.atomLevel = level
	})
}
