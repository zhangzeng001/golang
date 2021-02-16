# tailf模块使用

https://github.com/hpcloud/tail

https://pkg.go.dev/github.com/hpcloud/tail#section-documentation

https://www.cnblogs.com/wind-zhou/p/12840174.html

## tail包的作用

tail命令用途是依照要求将指定的文件的最后部分输出到标准设备，通常是终端，通俗讲来，就是把某个档案文件的最后几行显示到终端上，**假设该档案有更新，tail会自己主动刷新，确保你看到最新的档案内容** ，在日志收集中可以实时的监测日志的变化。

## 安装

```go
go get github.com/hpcloud/tail
```



## 用法

```go
package main

import (
	"fmt"
	"time"

	"github.com/hpcloud/tail"
)

// tailf的用法示例

func main() {
	fileName := "./my.log"
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return
	}
	var (
		line *tail.Line
		ok   bool
	)
	for {
		line, ok = <-tails.Lines //遍历chan，读取日志内容
		if !ok {
			fmt.Printf("tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("line:", line.Text)
	}
}
```



## 简单介绍

### type [Tail](https://github.com/hpcloud/tail/blob/master/tail.go#L76)

```go
type Tail struct {
    Filename string
    Lines    chan *Line
    Config

    tomb.Tomb // provides: Done, Kill, Dying
    // contains filtered or unexported fields
}
```



### func [TailFile](https://github.com/hpcloud/tail/blob/master/tail.go#L103)

```
func TailFile(filename string, config Config) (*Tail, error)
```

TailFile begins 传入参数：日志文件的路径和配置文件，返回一个指向Tail结构体对象的指针。

**config的数据结构为：**



### type [Config](https://github.com/hpcloud/tail/blob/master/tail.go#L58)

```go
type Config struct {
    // File-specifc
    Location    *SeekInfo // Seek to this location before tailing
    ReOpen      bool      // Reopen recreated files (tail -F)
    MustExist   bool      // Fail early if the file does not exist
    Poll        bool      // Poll for file changes instead of using inotify
    Pipe        bool      // Is a named pipe (mkfifo)
    RateLimiter *ratelimiter.LeakyBucket

    // Generic IO
    Follow      bool // Continue looking for new lines (tail -f)
    MaxLineSize int  // If non-zero, split longer lines into multiple lines

    // Logger, when nil, is set to tail.DefaultLogger
    // To disable logging: set field to tail.DiscardingLogger
    Logger logger
}
```

Config 用来定义文件被读取的方式。

Tail结构体中最重要的是Lines字段，他是存储Line指针来的一个通道。

Line的数据结构为：



### type [Line](https://github.com/hpcloud/tail/blob/master/tail.go#L28)

```
type Line struct {
    Text string
    Time time.Time
    Err  error // Error from tail
}
```

这个结构体是用来存储读取的信息。

> 最后简要总结一下流程：
>
> 1. 首先初始化配置结构体config
> 2. 调用TailFile函数，并传入文件路径和config，返回有个tail的结构体，tail结构体的Lines字段封装了拿到的信息
> 3. 遍历tail.Lnes字段，取出信息（注意这里要循环的取，因为tail可以实现实时监控）

贴上实战代码：

```go
package taillog

import (
	"fmt"
	"github.com/hpcloud/tail"
)


var (
	tailObj *tail.Tail
	//LogChan chan string
)


func Init (filename string)(err error){
	config := tail.Config{
		// File-specifc
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件那个位置开始读
		ReOpen:    true,                                 //是否重新打开
		MustExist: false,                 // Fail early if the file does notexist
		Poll:      true,                  // Poll for file changes instead of using inotify
		Follow:    true,                  // Continue looking for new lines (tail -f)

	}

	tailObj,err = tail.TailFile(filename, config) //TailFile(filename, config)

	if err != nil {
		fmt.Println("tail file err=", err)
		return
	}

	return

}

func ReadChan()<-chan *tail.Line{

	return tailObj.Lines 

}
```











