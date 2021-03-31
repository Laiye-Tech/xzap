
**⚠ ⚠ ⚠️因为go中协程panic没法在其他协程捕获，所以异步模式下，默认会丢日志，保证不丢日志，请在协程入口添加如下代码**

```
defer xzap.Logger(nil).Sync()
```



# xzap
提供兼容beego日志格式的zap库

## xzap/v0.1.3 更新

### 添加注入msgID，方法
```
ctx = xzap.InjectCtx(ctx, xzap.WithMsgID(2))
```

### 添加临时打开debug日志方法
```
ctx = xzap.InjectCtx(ctx, xzap.WithDebug())
```
这个ctx获取的logger会打debug级别的日志

在grpc请求的header传入debug，也会打开这条请求的debug级别的日志

通过grpcgatway调用的，在http的header中传入Grpc-Metadata-debug，打开这条请求的debug级别的日志



## 下载

go get github.com/Laiye-Tech/xzap

## 兼容性

1. 日志格式

beego
```bazaar
2019/11/15 18:06:58.453 [I] [xzap_test.go:103]  x 729219 729219 729219 729219 729219 729219 729219 729219
```

xzap
```bazaar
2019/11/15 18:07:05.419	[I]	[xzap/xzap_test.go:71]	x	{"": 729219, "": 729219, "": 729219, "": 729219, "": 729219, "": 729219, "": 729219, "": 729219}
```

xzap日期，级别，代码行数格式与beego相同；代码行号比beego多了目录。内容为json编码

日期与级别等部分分割为tab，日志平台可以被正确解析

2. 文件命名，日志旋转

支持安大小，日，旋转日志。不支持安小时旋转。支持清理过期日志。

beego日志命名
```bazaar
bee_log.log
bee_log.2019-11-14.001.log
```
xzap旋转日志命名
```bazaar
log.log
log-2019-11-14T20-03-54.810.log
```

3. 其他默认为异步模式，在异步模式下一定在main中退出时调用`logger.Sync()`，否则可能丢日志。
4. 若要打开同步模式，使用：xzap.WithAsyncParams(xzap.AsyncParams{Async:true})

## 使用
```bazaar
ctx := context.WithValue(context.Background(), "msgID", 1)
ctx = opentracing.ContextWithSpan(ctx, &jaeger.Span{})

xzap.InitLog("./log.log")
xzap.Sugar(ctx).Infof("hi %v", "world")
xzap.Logger(ctx).Info("hi", zap.Any("test", 1))
```

输出

```bazaar
2019/11/14 22:27:10.693	[I]	[xzap/xzap_test.go:22]	hi world	{"traceID": "0", "msgID": 1}
2019/11/14 20:39:26.172	[I]	[xzap/xzap_test.go:25]	hi	{"traceID": "0", "msgID": 1, "test": 1}
```

在ctx里注入traceID和msgID会自动记录到日志里，希望支持其他可以参考代码扩展。

`xzap.InitLog("./log.log")`参数

1. 保存30天
2. 安天旋转日志
3. 安大小旋转日志，默认256M
4. 异步刷新日志，默认500ms或100条

修改参数

```
xzap.InitLog("./log.log", xzap.WithDaily(false))
```

## 迁移

1. 方法1

```bazaar
log=xzap.Beego()
xzap.Beego().Info("hi %v", "world")
xzap.Beego().InfoCtx(ctx, "hi %v", "world")
```

注：只支持部分方法。

2. 方法2

```bazaar
log=xzap.Sugar()

log.Info改为log.Infof
```

3. 方法3

```bazaar
xzap.BeegoCtx(ctx).Info()
```

这种方法支持在日志中注入msgID，traceID

## 设置

1. 改为同步，或修改异步flush间隔
```
xzap.InitLog("",xzap.WithAsyncParams())
```

2. 打开或关闭日志daily切割
```
xzap.InitLog("",xzap.WithDaily())
```



