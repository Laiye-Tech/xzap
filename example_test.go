package xzap_test

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Laiye-Tech/xzap"
	"go.uber.org/zap"
)

func ExampleInitLog() {
	xzap.InitLog("./xx.log")
	defer zap.L().Sync()

	go http.ListenAndServe("localhost:8080", nil)

	//默认是 INFO 级别，DEUBG 日志不输出
	xzap.Sugar(nil).Infof("info should log")
	xzap.Logger(nil).Debug("debug should not log")

	// below code is equal to :
	//	  curl -XPOST "http://localhost:8080/logging" -d "level=debug&ttl=1s"
	contentType := "application/x-www-form-urlencoded"
	param := url.Values{
		"level": []string{"debug"},
		"ttl":   []string{"1s"},
	}
	resp, err := http.Post("http://localhost:8080/logging", contentType, strings.NewReader(param.Encode()))
	if err != nil {
		b, _ := ioutil.ReadAll(resp.Body)
		xzap.Sugar(nil).Fatalf("failed to request logging endpoint, err: %v, message: %v", err, string(b))
	}

	//DEUBG 日志会输出
	xzap.Sugar(nil).Infof("info should log")
	xzap.Logger(nil).Debug("debug should log")

	time.Sleep(time.Second + time.Millisecond)
	//回复到 INFO 级别，DEUBG 日志不输出
	xzap.Sugar(nil).Infof("info should log")
	xzap.Logger(nil).Debug("debug should not log")
}
