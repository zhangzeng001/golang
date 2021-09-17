# zap日志库使用

[github](https://github.com/uber-go/zap)

[doc](https://pkg.go.dev/go.uber.org/zap)

## 安装

运行下面的命令安装zap

```bash
go get -u go.uber.org/zap
```

## 配置Zap Logger

Zap提供了两种类型的日志记录器—`Sugared Logger`和`Logger`。

在性能很好但不是很关键的上下文中，使用`SugaredLogger`。它比其他结构化日志记录包快4-10倍，并且支持结构化和printf风格的日志记录。

在每一微秒和每一次内存分配都很重要的上下文中，使用`Logger`。它甚至比`SugaredLogger`更快，内存分配次数也更少，但它只支持强类型的结构化日志记录。

## Logger

```go
func NewProduction(options ...Option) (*Logger, error)
```

例子：

```go
package main

import (
	"go.uber.org/zap"
	"time"
)

func main() {
	// 获取logger对象
	logger,_ := zap.NewProduction()
	// 关闭logger对象
	defer logger.Sync()
	url := "http://www.baidu.com"
	logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}
```

结果：

```go
{"level":"info","ts":1613368072.9613779,"caller":"zap_logger/main.go:14","msg":"failed to fetch URL","url":"http://www.baidu.com","attempt":3,"backoff":1}
```

- 通过调用`zap.NewProduction()`/`zap.NewDevelopment()`或者`zap.Example()`创建一个Logger。
- 上面的每一个函数都将创建一个logger。唯一的区别在于它将记录的信息不同。例如production logger默认记录调用函数信息、日期和时间等。
- 通过Logger调用Info/Error等。
- 默认情况下日志都会打印到应用程序的console界面。

 日志记录器方法的语法是这样的： 

```go
func (log *Logger) MethodXXX(msg string, fields ...Field)
```

其中`MethodXXX`是一个可变参数函数，可以是Info / Error/ Debug / Panic等。每个方法都接受一个消息字符串和任意数量的`zapcore.Field`场参数。

每个`zapcore.Field`其实就是一组键值对参数。

```go
package main

import (
	"go.uber.org/zap"
	"net/http"
)

// 定义一个 *zap.logger类型的变量
var logger *zap.Logger

// 初始化logger
func InitLogger() {
	logger, _ = zap.NewProduction()
}

// 测试logger
func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	// func (log *Logger) Error|Info|...(msg string, fields ...Field)
	if err != nil {
		logger.Error(
			"Error fetching url..",
			zap.String("url", url),
			zap.Error(err))
	} else {
		logger.Info("Success..",
			zap.String("statusCode", resp.Status),
			zap.String("url", url),
			zap.String("test", "mytest"))
		resp.Body.Close()
	}
}

func main() {
	InitLogger()
	defer logger.Sync()
	simpleHttpGet("www.baidu.com")
	simpleHttpGet("http://www.baidu.com")
}
```

执行结果：

```go
{"level":"error","ts":1613367452.2733693,"caller":"zap_logger/main.go:21","msg":"Error fetching url..","url":"www.baidu.com","error":"Get \"www.baidu.com\": unsuppor
ted protocol scheme \"\"","stacktrace":"main.simpleHttpGet\n\tD:/go_code/src/zap_logger/main.go:21\nmain.main\n\tD:/go_code/src/zap_logger/main.go:37\nruntime.main\n
\tD:/Go/src/runtime/proc.go:204"}
{"level":"info","ts":1613367452.3754978,"caller":"zap_logger/main.go:26","msg":"Success..","statusCode":"200 OK","url":"http://www.baidu.com","test":"mytest"}
```



##  SugaredLogger

现在让我们使用Sugared Logger来实现相同的功能。

- 大部分的实现基本都相同。
- 惟一的区别是，我们通过调用主logger的`.Sugar()`方法来获取一个`SugaredLogger`。
- 然后使用`SugaredLogger`以`printf`格式记录语句

在性能不错但不是很关键的情况下，请使用`SugaredLogger`。它比其他结构化日志记录包快4-10倍，并且支持结构化和`printf`样式的日志记录 

```go
package main

import (
	"go.uber.org/zap"
	"time"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	url := "http://www.baidu.com"
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
}

/*
{"level":"info","ts":1613368356.315027,"caller":"zap_logger/main.go:13","msg":"failed to fetch URL","url":"http://www.baidu.com","attempt":3,"backoff":1}
{"level":"info","ts":1613368356.315027,"caller":"zap_logger/main.go:19","msg":"Failed to fetch URL: http://www.baidu.com"}
*/
```

```go
package main

import (
	"go.uber.org/zap"
	"net/http"
)

var sugarLogger *zap.SugaredLogger

func InitLogger() {
	logger, _ := zap.NewProduction()
	sugarLogger = logger.Sugar()
}

func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}

func main() {
	InitLogger()
	defer sugarLogger.Sync()
	simpleHttpGet("www.baidu.com")
	simpleHttpGet("http://www.baidu.com")
}

// 结果：
/*
D:\go_code\src\zap_logger>zap_logger.exe
{"level":"error","ts":1613368813.0136042,"caller":"zap_logger/main.go:19","msg":"Error fetching URL www.baidu.com : Error = Get \"www.baidu.com\": unsupported protoc
ol scheme \"\"","stacktrace":"main.simpleHttpGet\n\tD:/go_code/src/zap_logger/main.go:19\nmain.main\n\tD:/go_code/src/zap_logger/main.go:29\nruntime.main\n\tD:/Go/sr
c/runtime/proc.go:204"}
{"level":"info","ts":1613368813.103519,"caller":"zap_logger/main.go:21","msg":"Success! statusCode = 200 OK for URL http://www.baidu.com"}
*/
```



## 定制logger

### 日志文件写入文件

我们要做的第一个更改是把日志写入文件，而不是打印到应用程序控制台。

- 我们将使用`zap.New(…)`方法来手动传递所有配置，而不是使用像`zap.NewProduction()`这样的预置方法来创建logger。

```go
func New(core zapcore.Core, options ...Option) *Logger
```

 `zapcore.Core`需要三个配置——`Encoder`，`WriteSyncer`，`LogLevel`。 

 1.**Encoder**:编码器(如何写入日志)。我们将使用开箱即用的`NewJSONEncoder()`，并使用预先设置的`ProductionEncoderConfig()`。 

2.**WriterSyncer** ：指定日志将写到哪里去。我们使用`zapcore.AddSync()`函数并且将打开的文件句柄传进去。

```go
   file, _ := os.Create("./test.log")
   writeSyncer := zapcore.AddSync(file)
```

 3.**Log Level**：哪种级别的日志将被写入。 

我们将修改上述部分中的Logger代码，并重写`InitLogger()`方法。其余方法—`main()` /`SimpleHttpGet()`保持不变。 

```go
package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
)
// 定义全局*zap.SugaredLogger类型 sugarLogger 变量
var sugarLogger *zap.SugaredLogger

// getEncoder 构建NewJSONEncoder
func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}
// getLogWriter 构建 writeSyncer
func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create("./test.log")
	return zapcore.AddSync(file)
}
// InitLogger 初始化日志对象
func InitLogger() {
	// 构建New(core)
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core)
	// 获取sugar对象
	sugarLogger = logger.Sugar()
}

func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}

func main() {
	InitLogger()
	defer sugarLogger.Sync()
	simpleHttpGet("www.baidu.com")
	simpleHttpGet("http://www.baidu.com")
}
```

 当使用这些修改过的logger配置调用上述部分的`main()`函数时，以下输出将打印在文件——`test.log`中。 

### 将JSON Encoder更改为普通的Log Encoder

现在，我们希望将编码器从JSON Encoder更改为普通Encoder。为此，我们需要将`NewJSONEncoder()`更改为`NewConsoleEncoder()`。

```go
// getEncoder 构建NewJSONEncoder
func getEncoder() zapcore.Encoder {
	// 修改打印日志为非json
	return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
}
```

结果：

```go
1.6133699490115595e+09	debug	Trying to hit GET request for www.baidu.com
1.6133699490115595e+09	error	Error fetching URL www.baidu.com : Error = Get "www.baidu.com": unsupported protocol scheme ""
1.6133699490115595e+09	debug	Trying to hit GET request for http://www.baidu.com
1.6133699491022434e+09	info	Success! statusCode = 200 OK for URL http://www.baidu.com
```

### 更改时间编码并添加调用者详细信息

鉴于我们对配置所做的更改，有下面两个问题：

- 时间是以非人类可读的方式展示，例如1.572161051846623e+09
- 调用方函数的详细信息没有显示在日志中

我们要做的第一件事是覆盖默认的`ProductionConfig()`，并进行以下更改:

- 修改时间编码器
- 在日志文件中使用大写字母记录日志级别

```go
func getEncoder() zapcore.Encoder {
    // 非json输出
	encoderConfig := zap.NewProductionEncoderConfig()
    // zap.NewProductionEncoderConfig() json
    // 时间编码
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    // json类型时间字段key
    // MessageKey：输入信息的key名
    // LevelKey：输出日志级别的key名
    // TimeKey：输出时间的key名
    // encoderConfig.TimeKey = "time"
    // NameKey CallerKey StacktraceKey跟以上类似，看名字就知道
    // 一般zapcore.SecondsDurationEncoder,执行消耗的时间转化成浮点型的秒
    //  encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
    // EncodeCaller：一般zapcore.ShortCallerEncoder，以包/文件:行号 格式化调用堆栈
    // LineEnding：每行的分隔符。基本zapcore.DefaultLineEnding 即"\n"
    // 日志级别大写
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
```

 接下来，修改zap logger代码，添加将调用函数信息记录到日志中的功能。为此，我们将在`zap.New(..)`函数中添加一个`Option`。 

```go
logger := zap.New(core, zap.AddCaller())
// zap_logger/main.go:38 增加这一类信息
//2021-02-15T14:29:16.817+0800	DEBUG	zap_logger/main.go:38	Trying to hit GET request for www.baidu.com
```



```go
package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
)
// 定义全局*zap.SugaredLogger类型 sugarLogger 变量
var sugarLogger *zap.SugaredLogger

// getEncoder 构建NewJSONEncoder
func getEncoder() zapcore.Encoder {
	// 非json输出
	encoderConfig := zap.NewProductionEncoderConfig()
	// 时间编码
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 日志级别大写
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
	//return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
}
// getLogWriter 构建 writeSyncer
func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create("./test.log")
	return zapcore.AddSync(file)
}
// InitLogger 初始化日志对象
func InitLogger() {
	// 构建New(core)
	writeSyncer := getLogWriter()
	encoder := getEncoder()
    // zapcore.DebugLevel要输出的日志级别
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	// zap.AddCaller() 打印函数
	logger := zap.New(core,zap.AddCaller())
	// 获取sugar对象
	sugarLogger = logger.Sugar()
}

func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}

func main() {
	InitLogger()
	defer sugarLogger.Sync()
	simpleHttpGet("www.baidu.com")
	simpleHttpGet("http://www.baidu.com")
}
```

当使用这些修改过的logger配置调用上述部分的`main()`函数时，以下输出将打印在文件——`test.log`中。 

```go
2021-02-15T14:38:32.785+0800	debug	zap_logger/main.go:40	Trying to hit GET request for www.baidu.com
2021-02-15T14:38:32.802+0800	error	zap_logger/main.go:43	Error fetching URL www.baidu.com : Error = Get "www.baidu.com": unsupported protocol scheme ""
2021-02-15T14:38:32.802+0800	debug	zap_logger/main.go:40	Trying to hit GET request for http://www.baidu.com
2021-02-15T14:38:32.891+0800	info	zap_logger/main.go:45	Success! statusCode = 200 OK for URL http://www.baidu.com
```

### 修改时间格式

```go
EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
```

```go
func getEncoder() zapcore.Encoder {
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
```



## 增加自定义字段和值

### **with方式**

```
func (s *SugaredLogger) With(args ...interface{}) *SugaredLogger
```

```go
...
// InitLogger 初始化日志对象
func (z *zaplogger)InitLogger() *zap.SugaredLogger {
	// 构建New(core)
	writeSyncer := getLogWriter(z.logpath,z.isConsole)
	encoder := getEncoder(z.isConsole)
	core := zapcore.NewCore(encoder, writeSyncer, z.level)
	//core := zapcore.NewCore(encoder,  zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), z.level)
	// zap.AddCaller() 打印函数
	logger := zap.New(core,zap.AddCaller())
	// 获取sugar对象
	sugarLogger := logger.Sugar()
	return sugarLogger
}
...

logger := logobj.InitLogger()
logger2 := logger.With( "hello", "world",
			    "failure", errors.New("oh no"),
			    "count", 42,
			    "user", "zhang3")
		logger2.Info("xxxx")
```

结果：

```json
{"level":"INFO","ts":"2021-02-25 16:06:01","caller":"log_agent/main.go:46","msg":"xxxx","hello":"world","failure":"oh no","count":42,"user":"zhang3"}
```



### infow方式

```go
sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
```

```
{"level":"INFO","ts":"2021-02-25 16:13:20","caller":"log_agent/main.go:41","msg":"Infow自定义字段日志","name":"zhang3","age":12}
```



## 使用Lumberjack进行日志切割归档

使用第三方库[Lumberjack](https://github.com/natefinch/lumberjack)来实现。  日志切割归档功能

### 安装

执行下面的命令安装Lumberjack

```bash
go get -u github.com/natefinch/lumberjack
```

### zap logger中加入Lumberjack

要在zap中加入Lumberjack支持，我们需要修改`WriteSyncer`代码。我们将按照下面的代码修改`getLogWriter()`函数：

Lumberjack Logger采用以下属性作为输入:

- Filename: 日志文件的位置
- MaxSize：在进行切割之前，日志文件的最大大小（以MB为单位）
- MaxBackups：保留旧文件的最大个数，0无限制，与maxAge互相影响
- MaxAges：保留旧文件的最大天数
- Compress：是否压缩/归档旧文件

```go
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log", // 日志文件路径
		MaxSize:    10,           // 每个日志文件保存的大小 单位:M
		MaxBackups: 5,            // 日志文件最多保存多少个备份,0无限制，与maxAge互相影响
		MaxAge:     30,           // 日志文件最多保存多少天
		Compress:   false,        // 是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}
```

### 测试所有功能

 最终，使用Zap/Lumberjack logger的完整示例代码如下： 

```go
package main

import (
	"net/http"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugarLogger *zap.SugaredLogger

func main() {
	InitLogger()
	defer sugarLogger.Sync()
	simpleHttpGet("www.baidu.com")
	simpleHttpGet("http://www.baidu.com")
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
    // zapcore.DebugLevel要输出的日志级别
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 切割日志
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}
```

压缩日志文件

```go
test-2021-02-15T07-44-28.509.log
```

同时，可以在`main`函数中循环记录日志，测试日志文件是否会自动切割和归档（日志文件每1MB会切割并且在当前目录下最多保存5个备份）。

至此，我们总结了如何将Zap日志程序集成到Go应用程序项目中。

https://zhuanlan.zhihu.com/p/88856378?utm_source=wechat_session

https://mp.weixin.qq.com/s/i0bMh_gLLrdnhAEWlF-xDw

https://www.jianshu.com/p/5561396e61cf

https://www.jb51.net/article/180794.htm

https://blog.csdn.net/qq_27068845/article/details/103480451

https://www.cnblogs.com/chaselogs/p/9964424.html

# logrus日志库

https://github.com/go-ini/ini

https://www.liwenzhou.com/posts/Go/go_logrus/





















