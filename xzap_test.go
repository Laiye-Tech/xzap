package xzap

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
)

func TestMain(m *testing.M) {
	InitLog("./xx.log")
	defer zap.L().Sync()

	go http.ListenAndServe("localhost:8080", nil)

	m.Run()
}

func TestNew(t *testing.T) {
	log := New("./xx.log")
	defer log.Sync()
	log.Info("test")
}

func TestInitLog(t *testing.T) {
	ctx := context.WithValue(context.Background(), "msgID", 1)
	ctx = opentracing.ContextWithSpan(ctx, &jaeger.Span{})
	Sugar(ctx).Infof("hi %v", "world")
	Logger(ctx).Info("hi", zap.Any("test", 1))
	Beego().Info("hi %v", "world")
	Beego().Debug("debug hi %v", "world")
	Beego().InfoCtx(ctx, "hi %v", "world")
	//time.Sleep(time.Second)
	nctx := InjectCtx(ctx, WithDebug(), WithMsgID(2))
	Logger(nctx).Debug("debug hi")
	Logger(ctx).Debug("debug hi")
	logger.Sync()
}

func TestRuntimeLevel(t *testing.T) {
	Sugar(nil).Infof("info should log")
	Logger(nil).Debug("debug should not log")
	contentType := "application/x-www-form-urlencoded"
	param := url.Values{
		"level": []string{"debug"},
	}
	resp, err := http.Post("http://localhost:8080/logging", contentType, strings.NewReader(param.Encode()))
	assert.Nil(t, err)
	b, _ := ioutil.ReadAll(resp.Body)
	assert.Empty(t, string(b))
	Sugar(nil).Infof("info should log")
	Logger(nil).Debug("debug should log")
	logger.Sync()
}

func TestRuntimeLevelRestore(t *testing.T) {
	Sugar(nil).Infof("info should log")
	Logger(nil).Debug("debug should not log")
	contentType := "application/x-www-form-urlencoded"
	param := url.Values{
		"level": []string{"debug"},
		"ttl":   []string{"1s"},
	}
	resp, err := http.Post("http://localhost:8080/logging", contentType, strings.NewReader(param.Encode()))
	assert.Nil(t, err)
	b, _ := ioutil.ReadAll(resp.Body)
	assert.Empty(t, string(b))
	Sugar(nil).Infof("info should log")
	Logger(nil).Debug("debug should log")

	time.Sleep(time.Second + time.Millisecond)
	Sugar(nil).Infof("info should log")
	Logger(nil).Debug("debug should not log")
	logger.Sync()
}

func TestRuntimeLevelMultiRestore(t *testing.T) {
	Sugar(nil).Infof("info should log")
	Logger(nil).Debug("debug should not log")
	contentType := "application/x-www-form-urlencoded"
	{
		param := url.Values{
			"level": []string{"debug"},
			"ttl":   []string{"1m"},
		}
		resp, err := http.Post("http://localhost:8080/logging", contentType, strings.NewReader(param.Encode()))
		assert.Nil(t, err)
		b, _ := ioutil.ReadAll(resp.Body)
		assert.Empty(t, string(b))
		Sugar(nil).Infof("info should log")
		Logger(nil).Debug("debug should log")
	}
	{
		param := url.Values{
			"level": []string{"debug"},
			"ttl":   []string{"1s"},
		}
		resp, err := http.Post("http://localhost:8080/logging", contentType, strings.NewReader(param.Encode()))
		assert.Nil(t, err)
		b, _ := ioutil.ReadAll(resp.Body)
		assert.Empty(t, string(b))
		Sugar(nil).Infof("info should log")
		Logger(nil).Debug("debug should log")
	}
	time.Sleep(time.Second * 2)
	Sugar(nil).Infof("info should log")
	Logger(nil).Debug("debug should not log")
	logger.Sync()
}

func BenchmarkBenchWrite(b *testing.B) {
	b.ReportAllocs()
	f, err := ioutil.TempFile("", "")
	if err != nil {
		b.Fatal(err)
	}
	p := buffer.NewPool()
	bf := p.Get()
	bf.Reset()
	v := [1000]byte{}
	for i := 0; i < b.N; i++ {
		bf := p.Get()
		bf.Reset()
		for i := 0; i < 1000; i++ {
			bf.Write(v[:])
		}
		f.Write(bf.Bytes())
		bf.Free()
	}
}

func BenchmarkWrite(b *testing.B) {
	b.ReportAllocs()
	f, err := ioutil.TempFile("", "")
	if err != nil {
		b.Fatal(err)
	}
	v := [1000]byte{}
	for i := 0; i < b.N; i++ {
		for i := 0; i < 1000; i++ {
			f.Write(v[:])
		}
	}
}

func BenchmarkLog(b *testing.B) {
	os.Remove("./log.log")
	InitLog("./log.log", WithDaily(false))
	defer logger.Sync()
	b.ReportAllocs()
	v := 729219
	for i := 0; i < b.N; i++ {
		logger.Info("x",
			zap.Int("", v), zap.Int("", v), zap.Int("", v), zap.Int("", v),
			zap.Int("", v), zap.Int("", v), zap.Int("", v), zap.Int("", v),
		)
	}
}

func BenchmarkLogSync(b *testing.B) {
	os.Remove("./log.log")
	InitLog("./log.log", WithAsyncParams(AsyncParams{
		Async:         false,
		FlushInternal: 0,
		FlushLine:     0,
	}))
	defer logger.Sync()
	b.ReportAllocs()
	v := 729219
	for i := 0; i < b.N; i++ {
		logger.Info("x",
			zap.Int("", v), zap.Int("", v), zap.Int("", v), zap.Int("", v),
			zap.Int("", v), zap.Int("", v), zap.Int("", v), zap.Int("", v),
		)
	}
}

func BenchmarkBeeGoLog(b *testing.B) {
	os.Remove("./bee_log.log")
	log := beeLog("./bee_log.log")
	defer log.Flush()
	b.ReportAllocs()
	v := 729219
	for i := 0; i < b.N; i++ {
		log.Info("x %d %d %d %d %d %d %d %d", v, v, v, v, v, v, v, v)
	}
}

func beeLog(f string) *logs.BeeLogger {
	log := logs.NewLogger(100)
	//debug := configs.AppDebug()
	if true {
		log.EnableFuncCallDepth(true)
		log.SetLogFuncCallDepth(2)
		log.EnableFuncCallDepth(true)
	}
	logLevel := logs.LevelInfo
	log_config := map[string]interface{}{
		"filename": f,
		"maxsize":  1 << 28,
		"daily":    true,
		"maxdays":  30,
		"level":    logLevel,
	}
	file_config, _ := json.Marshal(log_config)
	log.SetLogger("file", string(file_config))
	log.Async()
	return log
}
