# zap日志库使用

[github](https://github.com/uber-go/zap)

[doc](https://pkg.go.dev/go.uber.org/zap)

##  SugaredLogger 

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
```



##  Logger 

**可靠性更高**，但仅支持强类型的结构化日志记录 

```go
package main

import (
	"fmt"
	"go.uber.org/zap"
	"time"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	url := "http://www.baidu.com"
	logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
	now := time.Now()
	unix := now.Unix()
	fmt.Println(unix)
}
```













https://mp.weixin.qq.com/s/i0bMh_gLLrdnhAEWlF-xDw

https://www.jianshu.com/p/5561396e61cf

https://www.jb51.net/article/180794.htm

https://blog.csdn.net/qq_27068845/article/details/103480451

https://www.cnblogs.com/chaselogs/p/9964424.html