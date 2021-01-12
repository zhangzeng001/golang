# flag命令行参数解析

import "flag"`

flag包实现了命令行参数的解析。

要求：

使用`flag.String(), Bool(), Int()`等函数注册flag

注册后必须使用 `flag.Parse()`初始化

## 基本使用

```go
package main

import (
	"flag"
	"fmt"
)

func main() {
    // 传入string类型参数，参数名ip，默认值127.0.0.1，帮助信息“xxxx”
	var ip = flag.String("ip","127.0.0.1", "help message for flagname")
    // 来解析命令行参数写入注册的flag里
	flag.Parse()
	fmt.Println("ip: ",*ip)
}
```

使用方法

​		使用`-未定义选项`返回错误信息

​		使用`-参数名 值`

```go
myflag.exe -ip 10.0.0.11
```



## 通过变量调用

```go
func main() {
	//var ip = flag.String("ip","127.0.0.1", "help message for flagname")
	//flag.Parse()
	//fmt.Printf("type:%T\n",flagvar)
	//fmt.Println("ip: ",*ip)

	var flagvar int
	flag.IntVar(&flagvar, "age", 0, "help message for age")
	flag.Parse()
	fmt.Printf("type:%T  %v\n",flagvar,flagvar)
}

// myflag.exe -age 22
// 22
```



## 传入自定义类型

 或者你可以自定义一个用于flag的类型<font color=FF0000>**(满足Value接口)**</font>并将该类型用于flag解析，如下： 

```
flag.Var(&flagVal, "name", "help message for flagname")
```



源码解析：

​		参数值为value 传入参数；name参数名；usage帮助信息

​		其中value必须有String()和Set()两个方法

```go
func Var(value Value, name string, usage string) {
	CommandLine.Var(value, name, usage)
}

type Value interface {
	String() string
	Set(string) error
}
```



例子：

​		传入一个json的参数，并解析出相关字段（json中双引号需转义，不然传入后字段和值双引号会消失无法解析）

```go
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
)

// Students 定义一个结构体用于接收传入参数，字段需相同
type Students struct {
	Name string
	Age int
}

// String 构建必须的String()方法，按要求返回字符串
func (s *Students)String() string {
	//fmt.Println("string函数-----------------",fmt.Sprintf("%v", *s))
	return fmt.Sprintf("%v", *s)
}

// Set 构建必须的set方法，返回error；其中参数v是传入的json字符串，只需解析参数v 即可
func (s *Students)Set(v string) error{
    // 定义一个Students结构类型的变量接收遍历后的json，也可以直接传s
	var tmp Students
	// 解析json字符串
	err := json.Unmarshal([]byte(v),&tmp)
	if err != nil{
		fmt.Println("非Json类型参数",err)
		return err
	}
    // 判断传入json是否为空
	if tmp == (Students {}){
		return errors.New("非法参数")
	}
    // 赋值，其中s即使传入结构体jsonFlag
	*s = tmp
    // 无异常返回nil
	return nil
}

// 定义接收类型
var jsonFlag Students

// 放到初始化函数，被首先执行
func init()  {
	flag.Var(&jsonFlag, "loadjson","传入json类型参数")
}

func main() {
    // 传参初始化，写入parse
	flag.Parse()
    // 使用获取的参数值
	fmt.Println(jsonFlag.Name)
	fmt.Println(jsonFlag.Age)
}
```



## 传入多个参数



```go
func main() {
	var arg1 string
	var arg2 int
    
	flag.StringVar(&arg1,"name","xxx","Input name")
	flag.IntVar(&arg2,"age",0,"Input age")
    
	flag.Parse()
    
	fmt.Println(arg1)
	fmt.Println(arg2)
}
```

返回：

```go
zhang3
12
```



## 传参模拟copy命令

[文件操作]()



# 第三方库 pflag 



#  os.Args 







