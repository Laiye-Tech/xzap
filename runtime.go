package xzap

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	helpMessage     = "You can modify logging level by passing 'level' param with one of debug/info/warn/error. \nFor example:\n\tlevel=debug"
	errParseForm    = "failed to parse form from query."
	errInvalidLevel = "invalid logging level."

	defaultTTL = time.Minute * 15
	reqchan    = make(chan *modRtLevelRequest)
)

var requestLevel2zap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
}

type modRtLevelRequest struct {
	level zapcore.Level
	ttld  time.Duration
	done  chan struct{}
}

func addRuntimeLevelHandler() {
	http.HandleFunc("/logging", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			_, _ = w.Write([]byte(errParseForm))
			return
		}
		level := r.PostForm.Get("level")
		if level == "" {
			_, _ = w.Write([]byte(helpMessage))
			return
		}
		zapLevel, prs := requestLevel2zap[level]
		if !prs {
			_, _ = w.Write([]byte(errInvalidLevel))
			return
		}

		ttld := defaultTTL
		if ttl := r.Form.Get("ttl"); ttl != "" {
			d, err := time.ParseDuration(ttl)
			if err == nil {
				ttld = d
			}
		}

		modReq := &modRtLevelRequest{
			level: zapLevel,
			ttld:  ttld,
			done:  make(chan struct{}),
		}

		zap.S().Warnf("changing logging level to %v for %v", modReq.level, modReq.ttld)
		reqchan <- modReq
		<-modReq.done
		zap.S().Warnf("logging level is changed to %v for %v", modReq.level, modReq.ttld)
	})
}

//运行时替换全局logger
func startRuntimeLevelServer() {
	addRuntimeLevelHandler()

	var (
		tk      = &time.Ticker{} // 到期时回复日志到初始化级别，多次修改请求会覆盖恢复时机
		restore func()           // 恢复全局logger初始化时配置的日志级别
	)

	for {
		select {
		case req, more := <-reqchan:
			if !more {
				return
			}
			originLevel := atomLevel.Level()
			if restore == nil {
				restore = func() {
					atomLevel.SetLevel(originLevel)
					zap.S().Infof("logging level is restored to %v", originLevel)
				}
			}

			atomLevel.SetLevel(req.level)

			tk.Stop()
			tk = time.NewTicker(req.ttld)
			close(req.done)
		case <-tk.C:
			if restore != nil {
				restore()
				tk.Stop()
			}
		}
	}
}
