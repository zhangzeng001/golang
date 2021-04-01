# JSON 序列化

golang 导入json包支持

`import "encoding/json"`

json序列化方法

`json.Marshal(v interface{}) ([]byte,error)`

其中：

​	v 是需要转化为json的类型值

<font color=FF0000>**注意：**</font>Marshal()函数只有在转换成功的时候才会返回数据，在转换的过程中需要注意以下几点：

* JSON对象只支持string作为key，所以要编码一个map，必须是map[string]Type 这种类型
* struct key 需大写
* channel、complex、function是不能被编码成JSON的
* 指针在编码的时候会输出指针指向的内容，二空指针会输出null

## map 转 JSON

```go
package main

import (
	"encoding/json"
	"fmt"
)

// map 转 Json

func main(){
	s1 := make([]string,2,2)
	s1 = []string{"a","b"}
	// 定义一个map变量并初始化
	m := map[string][]string{
		"zhang3": s1,
		"li4": []string{"c","d"},
	}
	fmt.Println(m) // map[li4:[c d] zhang3:[a b]]
	// 将 map 解析为json格式
	if data,err := json.Marshal(m) ; err==nil{
		fmt.Printf("%s\n",data)  //{"li4":["c","d"],"zhang3":["a","b"]}
		fmt.Printf("%T\n",data)  //[]uint8
	}
}
```

`Marshal()`函数返回的JSON字符串是没有空白字符和缩进的，这种紧凑的表示形式是最常用的传输形式，但不好阅读，。如果需要为前端输出便于阅读，可以调用`json.MarshalIndent()`，该函数有两个参数表示每一行的前缀和缩进方式，如下：

`func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)`

```go
package main

import (
	"encoding/json"
	"fmt"
)

// map 转 Json

func main(){
	s1 := make([]string,2,2)
	s1 = []string{"a","b"}
	// 定义一个map变量并初始化
	m := map[string][]string{
		"zhang3": s1,
		"li4": []string{"c","d"},
	}
	fmt.Println(m) // map[li4:[c d] zhang3:[a b]]
	// 将 map 解析为json格式,MarshalIndent为重端便于阅读的格式
	if data,err := json.MarshalIndent(m,"","    ") ; err==nil{
		fmt.Printf("%s\n",data)
		fmt.Printf("%T\n",data)  //[]uint8
	}
}

// 返回
{
    "li4": [
        "c",
        "d"
    ],
    "zhang3": [
        "a",
        "b"
    ]
}
```



## 结构体 转 JSON

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Student struct {
	Name string   // 必须大写开头
	Course string
}

type Adderss struct {
	province string
}

func main(){
	p1 := Student{
		//id : 100,
		Name: "zhang3",
		Course: "GoLang",
	}

	if data,err := json.Marshal(p1) ; err==nil{
		fmt.Printf("json:%#s\n", data) //json:{"Name":"zhang3","Course":"GoLang"}
	}
}

```



```go
package main

import (
	"encoding/json"
	"fmt"
)

type DebugInfo struct {
	Level string
	Msg string
	author string
}

func main(){
	// 定义一个结构体
	dbgInfs := []DebugInfo{
		DebugInfo{Level: "debug",Msg: `File: "test.txt" Not Found`,author: "Cynhard"},
		DebugInfo{Level: "",Msg: "Logic error",author: "Gopher"},
	}
	// 将结构体转换成JSON格式
	if data,err := json.Marshal(dbgInfs);err == nil{
		fmt.Printf("%s",data)
	}
}
```



# JSON 解析

## JSON转切片类型的map

```go
package main

import (
	"encoding/json"
	"fmt"
)

func main(){
	// 定义json格式的字符串
	data := `[{"name":"zhang3","age":22},{"name":"li4","age":18}]`
	var dbgInfos []map[string]string
	// 将字符串解析成map切片
	json.Unmarshal([]byte(data),&dbgInfos)
	fmt.Println(dbgInfos)
}
```



## JSON转结构体

<font color=F0000>接受JSON的字段名必须大写，否则无法解析</font>

```go
package main

import (
	"encoding/json"
	"fmt"
)

type DebugInfo struct {
	Level string
	Msg string
	author string // 小写开头未导出字段不会被json解析
}


func main()  {
	data := `[{"Level": "debug","Msg": "File Not Found","author": "Cynhard"},{"Level": "","Msg": "Logic error","author": "Gopher"}]`
	var dbgInfos []DebugInfo
	err := json.Unmarshal([]byte(data),&dbgInfos)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(dbgInfos) // [{debug File Not Found } { Logic error }]
}
```



```go
type DebugInfo struct {
	Level string
	Msg string
	author string // 未导出字段不会被json解析
}


func main()  {
	data := `{"Level": "debug","Msg": "File Not Found","author": "Cynhard"}`
	var dbgInfos DebugInfo
	err := json.Unmarshal([]byte(data),&dbgInfos)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(dbgInfos) // [{debug File Not Found } { Logic error }]
}
```





















