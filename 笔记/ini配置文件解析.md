# ini配置文件解析

https://github.com/go-ini/inihttps://github.com/go-ini/ini

https://ini.unknwon.io/docs/intro/getting_started

## 安装

```
go get gopkg.in/ini.v1
```



## 使用

加载配置文件

```
func Load(source interface{}, others ...interface{}) (*File, error)
```



示例ini文件

```ini
# possible values : production, development
app_mode = development

[paths]
# Path to where grafana can store temp files, sessions, and the sqlite3 db (if that is used)
data = /home/git/grafana
logs = /home/logs/grafana.log

[server]
# Protocol (http or https)
protocol = http

# The http port  to use
http_port = 9999

# Redirect to correct domain if host header does not match domain
# Prevents DNS rebinding attacks
enforce_domain = true
```

```go
package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

func main() {
	// 加载配置文件
	cfg,err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	//// 获取分组下所有配置
	//paths := cfg.Section("paths")
	//fmt.Println(paths.Body())

	// 无分组配置获取
	app_mode := cfg.Section("").Key("app_mode")
	fmt.Println(app_mode)   // development
	// 获取分组下string配置
	logs := cfg.Section("paths").Key("logs").String()
	fmt.Println(logs)      // /home/logs/grafana.log

	// 获取分组下int配置,非法值转换为0
	port,_ := cfg.Section("server").Key("http_port").Int()
	fmt.Println(port)     // 9999
	// 试一试自动类型转换，非法值返回值MustInt()
	fmt.Printf("Port Number: (%[1]T) %[1]d\n", cfg.Section("server").Key("http_port").MustInt(9999))
	fmt.Printf("Enforce Domain: (%[1]T) %[1]v\n", cfg.Section("server").Key("enforce_domain").MustBool(false))
	// 如果读取的值不在候选列表内，则会回退使用提供的默认值
	fmt.Println("Email Protocol:",
		cfg.Section("server").Key("protocol").In("smtp", []string{"imap", "smtp"}))

	// 修改某个值然后进行保存为新文件
	cfg.Section("").Key("app_mode").SetValue("production")
	cfg.SaveTo("my.ini.local")
}
```

### 配置映射到结构体

https://ini.unknwon.io/docs/advanced/map_and_reflect

```ini
Name = Unknwon
age = 21
Male = true
Born = 1993-01-01T20:17:05Z

[Note]
Content = Hi is a good man!
Cities = HangZhou, Boston
```

```go
package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"time"
)

type Note struct {
	Content string
	Cities  []string
}

type Person struct {
	Name string
	Age  int `ini:"age"`
	Male bool
	Born time.Time
	Note
	Created time.Time `ini:"-"`
}

func main() {
	cfg, err := ini.Load("config.ini")
	// ...
	p := new(Person)
	err = cfg.MapTo(p)
	fmt.Println(err)
	// ...

	// 一切竟可以如此的简单。
	err = ini.MapTo(p, "config.ini")
	fmt.Println(p.Age)
	// ...

	// 嗯哼？只需要映射一个分区吗？
	n := new(Note)
	err = cfg.Section("Note").MapTo(n)
	fmt.Println(n)
	fmt.Println(n.Content)
	// ...
}
```

结构的字段怎么设置默认值呢？很简单，只要在映射之前对指定字段进行赋值就可以了。如果键未找到或者类型错误，该值不会发生改变。 

```go
// ...
p := &Person{
    Name: "Joe",
}
// ...// ...
p := &Person{
    Name: "Joe",
}
// ...
```



### 从结构反射

### 配合 ShadowLoad 进行映射

如果您希望配合 [ShadowLoad](https://ini.unknwon.io/docs/howto/work_with_keys#same-key-with-multiple-values) 将某个分区映射到结构体，则需要指定 `allowshadow` 标签。

假设您有以下配置文件：

```ini
[IP]
value = 192.168.31.201
value = 192.168.31.211
value = 192.168.31.221
```

您应当通过如下方式定义对应的结构体：

```go
type IP struct {
   Value    []string `ini:"value,omitempty,allowshadow"`
}
```

如果您不需要前两个标签规则，可以使用 `ini:",,allowshadow"` 进行简写。

