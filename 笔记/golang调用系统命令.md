# Golang 调用 Linux 命令

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

