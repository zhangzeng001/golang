# Golang 调用 Linux 命令

http://www.zzvips.com/article/63368.html

https://www.linuxprobe.com/golang-exec-command.html

**Golang** 中可以使用 `os/exec` 来执行 **Linux** 命令，下面是一个简的示例：

```go
package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func main() {
	cmd := exec.Command("/bin/bash", "-c", `df -lh`)

	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return
	}

	//读取所有输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
		return
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return
	}
	fmt.Printf("stdout:\n\n %s", bytes)
}
```

或者创建一个缓冲读取器按行读取：

```go
package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("/bin/bash", "-c", `df -lh`)

	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return
	}

	//使用带缓冲的读取器
	outputBuf := bufio.NewReader(stdout)

	for {

		//一次获取一行,_ 获取当前行是否被读完
		output, _, err := outputBuf.ReadLine()
		if err != nil {

			// 判断是否到文件的结尾了否则出错
			if err.Error() != "EOF" {
				fmt.Printf("Error :%s\n", err)
			}
			return
		}
		fmt.Printf("%s\n", string(output))
	}

	//wait 方法会一直阻塞到其所属的命令完全运行结束为止
	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return
	}
}
```

输出结果：

![img](http://cdn.tianfeiyu.com/pipe_command.png)

在写这句 `if err.Error() != "EOF"` 时，一直以为可以直接将 `error` 类型直接转为 `string` 然后就可以比较了，所以刚开始写的代码是这样的 `if string(err) != "EOF"`,但是一直报下面这个错误：

```go
\# command-line-arguments
./exec_command.go:36: cannot convert err (type error) to type string
```

于是查了下才明白，`error` 类型本身是一个预定义好的接口，里面定义了一个`method`：

```go
type error interface {
    Error() string
}
```

所以 `err.Error()` 才是一个 `string` 类型的返回值。



## 得到错误输出

```go
package main

import (
	"io/ioutil"
	"log"
	"os/exec"
)

//runInLinux 执行系统命令
func RunLinuxCmd(intputCmd string) (string, error) {
	//fmt.Println("Running Linux cmd:" + cmd)
	//创建获取命令输出管道
	cmd := exec.Command("/bin/bash", "-c", intputCmd)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return "",err
	}
    // --------------------这里----------------------
	stderr,_ := cmd.StderrPipe()

	//执行命令
	if err = cmd.Start(); err != nil {
		log.Println("Error:The command is err,", err)
		return "",err
	}

	//读取所有正常输出
	stdoutBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println("==> ",intputCmd," ReadAll Stdout:", err.Error())
		return "",err
	}
    //--------------------这里----------------------
	// 读取所有错误输出
	stderrBytes, err := ioutil.ReadAll(stderr)
	if err != nil {
		log.Println("==> ",intputCmd," ReadAll Stdout:", err.Error())
		return "",err
	}

	if err = cmd.Wait(); err != nil {
        // --------------------这里----------------------
		// log.Printf("==> %v wait: %v", intputCmd,string(stderrBytes))
		log.Println("==> 命令执行失败",intputCmd," wait:", err.Error())
		return string(stderrBytes),err
	}
	log.Printf("==> 命令执行成功: %s\n", intputCmd)
	return string(stdoutBytes), err

	//result, err := exec.Command("/bin/sh", "-c", intputCmd).Output()
	//if err != nil {
	//	return "", err
	//}
	//return strings.TrimSpace(string(result)), err
}

func main() {
    res,_ := RunLinuxCmd("ifconfig")
    fmt.Println(res)
}
```





linux打包

```go
set GOARCH=amd64

set GOOS=linux

go build main.go

此时会生成一个没有后缀的二进制文件  main
```



交互命令

https://blog.csdn.net/dirk2014/article/details/53700435









# 示例：

## 通过接口重启某个只能单次执行的服务

```go
package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

var s = &http.Server{
	Addr:           ":8080",
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
}


func execCmd()(msg string,err error){
	cmd := exec.Command("/bin/sh", "-c", "nginx -s reload")
	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return "",err
	}

	//执行命令
	if err = cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return"",err
	}

	//使用带缓冲的读取器
	reader := bufio.NewReader(stdout)
	for {
		//一次获取一行,_ 获取当前行是否被读完
		line, _, err := reader.ReadLine()
		if err != nil {
			// 判断是否到文件的结尾了否则出错
			if err.Error() != "EOF" {
				fmt.Printf("Error :%s\n", err)
			}
			if len(line) !=0 {
				msg = msg + string(line) // 最后一行
			}
			break
		}
		msg = msg+string(line)
		fmt.Println(msg)
	}

	//wait 方法会一直阻塞到其所属的命令完全运行结束为止
	err = cmd.Wait()
	if err != nil{
		fmt.Println(err)
		return "",err
	}
	fmt.Println(msg)
	return msg,nil
}


func restart (w http.ResponseWriter,r *http.Request){
	//fmt.Fprintln(w,"<h1>restart</h1>")
	//msg, _ := json.Marshal(`{Code: 200, Msg: "验证成功"}`)
	msg,err := execCmd()
	if err != nil{
		//fmt.Sprintln(err)
		w.Write([]byte(fmt.Sprintln(err)))
		return
	}
	w.Write([]byte(msg))
}

func main() {
	http.HandleFunc("/",restart)
	err := s.ListenAndServe()
	if err != nil{
		fmt.Println(err)
		return
	}
}
// 启动一个socket 端口监听 8080

// http 接收get请求并执行程序重启操作后返回执行结果
```





## 通过脚本，执行一次就阻塞

运维文件---golang



# [go语言中os/signal包的学习与使用](https://www.cnblogs.com/smallleiit/p/10844728.html)

```go
package main;
 
import (
    "os"
    "os/signal"
    "fmt"
)
 
//signal包中提供了两个函数
//Notifyf()用于监听信号
//Stop()用于停止监听
 
func main()  {
    ch := make(chan os.Signal);
    //notify用于监听信号
    //参数1表示接收信号的channel
    //参数2及后面的表示要监听的信号
    //os.Interrupt 表示中断
    //os.Kill 杀死退出进程
    signal.Notify(ch, os.Interrupt, os.Kill);
 
    //获取信号，如果没有会一直阻塞在这里。
    s := <-ch;
    //我们通过Ctrl+C或用taskkill /pid -t -f来杀死进程，查看效果。
    fmt.Println("信号：", s);
}<br><br><br>
```

```go
package main
 
import (
    "os"
    "os/signal"
    "fmt"
)
 
func main() {
    ch := make(chan os.Signal);
    //如果不指定要监听的信号，那么默认是所有信号
    signal.Notify(ch);
 
    //停止向ch转发信号，ch将不再收到任何信号
    signal.Stop(ch);
    //ch将一直阻塞在这里，因为它将收不到任何信号
    //所以下面的exit输出也无法执行
    <-ch;
    fmt.Println("exit");
}
```

## go语言中signal.Notify

https://blog.csdn.net/weixin_39172380/article/details/103408346

## go中的信号量

| 信号    | 值       | 动作 | 说明                                                         |
| ------- | -------- | ---- | ------------------------------------------------------------ |
| SIGHUP  | 1        | Term | 终端控制进程结束(终端连接断开)                               |
| SIGINT  | 2        | Term | 用户发送INTR字符(Ctrl+C)触发                                 |
| SIGQUIT | 3        | Core | 用户发送QUIT字符(Ctrl+/)触发                                 |
| SIGILL  | 4        | Core | 非法指令(程序错误、试图执行数据段、栈溢出等)                 |
| SIGABRT | 6        | Core | 调用abort函数触发                                            |
| SIGFPE  | 8        | Core | 算术运行错误(浮点运算错误、除数为零等)                       |
| SIGKILL | 9        | Term | 无条件结束程序(不能被捕获、阻塞或忽略)                       |
| SIGSEGV | 11       | Core | 无效内存引用(试图访问不属于自己的内存空间、对只读内存空间进行写操作) |
| SIGPIPE | 13       | Term | 消息管道损坏(FIFO/Socket通信时，管道未打开而进行写操作)      |
| SIGALRM | 14       | Term | 时钟定时信号                                                 |
| SIGTERM | 15       | Term | 结束程序(可以被捕获、阻塞或忽略)                             |
| SIGUSR1 | 30,10,16 | Term | 用户保留                                                     |
| SIGUSR2 | 31,12,17 | Term | 用户保留                                                     |
| SIGCHLD | 20,17,18 | Ign  | 子进程结束(由父进程接收)                                     |
| SIGCONT | 19,18,25 | Cont | 继续执行已经停止的进程(不能被阻塞)                           |
| SIGSTOP | 17,19,23 | Stop | 停止进程(不能被捕获、阻塞或忽略) SIGTSTP 18,20,24 Stop 停止进程(可以被捕获、阻塞或忽略) SIGTTIN 21,21,26 Stop 后台程序从终端中读取数据时触发 SIGTTOU 22,22,27 Stop 后台程序向终端中写数据时触发 |

*有些信号名对应着3个信号值，这是因为这些信号值与平台相关*
*SIGKILL和SIGSTOP这两个信号既不能被应用程序捕获，也不能被操作系统阻塞或忽略*

## kill与kill9的区别

- kill pid的作用是向进程号为pid的进程发送SIGTERM（这是kill默认发送的信号），该信号是一个结束进程的信号且可以被应用程序捕获。若应用程序没有捕获并响应该信号的逻辑代码，则该信号的默认动作是kill掉进程。这是终止指定进程的推荐做法。
- kill -9 pid则是向进程号为pid的进程发送SIGKILL（该信号的编号为9），从本文上面的说明可知，SIGKILL既不能被应用程序捕获，也不能被阻塞或忽略，其动作是立即结束指定进程。通俗地说，应用程序根本无法“感知”SIGKILL信号，它在完全无准备的情况下，就被收到SIGKILL信号的操作系统给干掉了，显然，在这种“暴力”情况下，应用程序完全没有释放当前占用资源的机会。事实上，SIGKILL信号是直接发给init进程的，它收到该信号后，负责终止pid指定的进程。在某些情况下（如进程已经hang死，无法响应正常信号），就可以使用kill -9来结束进程。
- 若通过kill结束的进程是一个创建过子进程的父进程，则其子进程就会成为孤儿进程（Orphan Process），这种情况下，子进程的退出状态就不能再被应用进程捕获（因为作为父进程的应用程序已经不存在了），不过应该不会对整个linux系统产生什么不利影响。

## 优雅地退出程序

在长时间的程序运行过程中，可能有大量的系统资源被申请，无论在以何种方式退出前，他们应该及时将这些资源释放并将状态输出到日志中方便调试和排错。

signal.Notify方法监听和捕获信号量

> func Notify(c chan<- os.Signal, sig …os.Signal)

首先定义一个chan传递信号量，然后说明那些信号量是需要被捕获的（不填的话就默认捕获任何信号量）

```go
sc := make(chan os.Signal, 1)
signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
12
```

监听指定信号量

```go
EXIT:
	for {
		sig := <-sc
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Panic("SIGQUIT")
		case syscall.SIGHUP:
			log.Panic("SIGHUP")
		case syscall.SIGHUP:
			log.Panic("SIGINT")
		default:
			break EXIT
		}
	}
```



# 禁止 main 函数退出的方法

```go

func main() {
	defer func() { for {} }()
}
 
func main() {
	defer func() { select {} }()
}
 
func main() {
	defer func() { <-make(chan bool) }()
}

```







# ssh连接执行交互命令

https://www.cnblogs.com/zhzhlong/p/12552410.html









