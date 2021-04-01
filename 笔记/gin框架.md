# 参考

[官方文档](https://gin-gonic.com/zh-cn/docs/)

[github](https://github.com/gin-gonic/gin#installation)

# 安装

1.下载并安装 gin：

```shell
$ go get -u github.com/gin-gonic/gin
```

2.将 gin 引入到代码中：

```go
import "github.com/gin-gonic/gin"
```

3.（可选）如果使用诸如 `http.StatusOK` 之类的常量，则需要引入 `net/http` 包：

```go
import "net/http"
```

4.启动项目

```sh
$ go run main.go
```



# resetfull 接口

```
func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.GET("/someGet", getting)
	router.POST("/somePost", posting)
	router.PUT("/somePut", putting)
	router.DELETE("/someDelete", deleting)
	router.PATCH("/somePatch", patching)
	router.HEAD("/someHead", head)
	router.OPTIONS("/someOptions", options)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}
```



# AsciiJSON

 使用 AsciiJSON 生成具有转义的非 ASCII 字符的 ASCII-only JSON。 

```golang
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	// 定义一个ping接口，返回pong
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
/* ####################################################################### */
	type d struct {
		S1 interface{} `json:"s1"`
		S2 interface{} `json:"ssss2"`
	}

	// 构建一个数据作为返回json
	data1 := d{
		S1: "111",
		S2: "222",
	}

	//data := map[string]interface{}{
	//	"lang": "GO语言",
	//	"tag":  "<br>",
	//}

	r.GET("/someJSON", func(c *gin.Context) {
		// 输出 : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
		//c.AsciiJSON(http.StatusOK, data)
		
		// 输出 : {"s1": "111","ssss2": "222"}
		c.AsciiJSON(http.StatusOK, data1) 
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

```



# 路由

## 路由参数

单匹配

```shell
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 定义一个处理路由的函数，传参必须为(c *gin.Context)，Param为取值对应路由定义:xxx
func users(c *gin.Context){
	name := c.Param("name")
	c.String(http.StatusOK,"Hello %s",name)
}

func main() {
	r := gin.Default()
	// 定义一个路由，定义获取参数key为name，param取值的名称
	r.GET("/user/:name" ,users)

	r.Run(":8080")
}
```

结果：

```http
# GET http://localhost:8080/user/zhang3
Hello zhang3
```



匹配多个

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 定义一个处理路由的函数，传参必须为(c *gin.Context)，Param为取值对应路由定义:xxx
func users(c *gin.Context){
	action := c.Param("acion")
	//name := c.Param("name")
	//message := name + action
	//c.String(http.StatusOK,"Hello %s",message)
	c.String(http.StatusOK,"Hello %s",action)
}

func main() {
	r := gin.Default()
	// 定义一个路由，定义匹配的参数key为action，param取值的名称
	r.GET("/user/:name/*action" ,users)

	r.Run(":8080")
}
```

结果：

```http
# GET http://localhost:8080/user/zhang3/asdfasdfaction/adsfasdf
Hello /asdfasdfaction/adsfasdf
```



## url参数

```go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func main() {
	r := gin.Default()
	r.POST("/login",login)

	r.Run(":8080")
}

func login(c *gin.Context)  {
	var user User
	err := c.ShouldBindQuery(&user)
	if err != nil {
		fmt.Println("参数错误")
		c.String(http.StatusOK,"参数错误!")
	}
	message := user.Username +" "+ user.Password
	c.String(http.StatusOK,"Url Params is %s",message)
}
```

结果

```shell
$ POST http://localhost:8080/login?username=zhang3&password=123
Url Params is zhang3 123
```



## body参数

* 简单示例

  ```go
  package main
  
  import (
  	"github.com/gin-gonic/gin"
  	"net/http"
  )
  
  func main() {
  	router := gin.Default()
  	router.POST("/login", events)
  	router.Run(":8080")
  }
  
  func events(c *gin.Context) {
  	buf := make([]byte, 1024)
  	n, _ := c.Request.Body.Read(buf)
  	c.String(http.StatusOK,string(buf[0:n]))
  }
  ```

  执行结果

  ```shell
  $ POST http://localhost:8080/login?username=zhang3&password=123
  {
      "username":"zhang3",
      "password":"123"
  }
  ```

  

* body转结构体

  ```go
  package main
  
  import (
  	"encoding/json"
  	"github.com/gin-gonic/gin"
  	"net/http"
  )
  
  // 接收body json的结构体
  type User struct {
  	Username string `form:"username"`
  	Password string `form:"password"`
  }
  
  func main() {
  	router := gin.Default()
  	router.POST("/login", events)
  	router.Run(":8080")
  }
  
  func events(c *gin.Context) {
      // 初始化结构体
  	var user User
      // 接受数据
  	buf := make([]byte, 1024)
  	n, _ := c.Request.Body.Read(buf)
      // json解析到结构体
      _ = json.Unmarshal([]byte((buf[0:n])),&user)
  
  	c.String(http.StatusOK,user.Username+" "+user.Password)
  }
  ```

  执行结果：

  ```shell
  $ POST http://localhost:8080/login
  zhang3 123
  ```

  





## 路由组

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	v1 := r.Group("v1")
	{
		v1.GET("/login",login)
		v1.GET("/status",status)
	}
	v2 := r.Group("v2")
	{
		//v2.GET("/xxxx",xxxx)
	}
	r.Run(":8080")
}
func login(c *gin.Context)  {
	c.String(http.StatusOK,"login")
}
func status(c *gin.Context)  {
	c.String(http.StatusOK,"status")
}
```

结果：

```shell
$ GET http://localhost:8080/v1/login
login
$ GET http://localhost:8080/v1/status
status
```









# HTML 渲染

 使用 `LoadHTMLGlob() `或者 `LoadHTMLFiles() `

* `LoadHTMLGlob()`

  ```go
  package main
  
  import (
  	"github.com/gin-gonic/gin"
  	"net/http"
  )
  
  func main() {
  	router := gin.Default()
  	router.LoadHTMLGlob("templates/*")
  	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
  	router.GET("/index", func(c *gin.Context) {
  		c.HTML(http.StatusOK, "index.tmpl", gin.H{
  			"title": "Index Body ...",
  		})
  	})
  	router.Run(":8080")
  }
  ```

   templates/index.tmpl 

  ```html
  <html>
  	<h1>
  		{{ .title }}
  	</h1>
  </html>
  ```

  **使用不同目录下名称相同的模板**

  `router.LoadHTMLGlob("templates/**/*")`使用`{{ define "posts/index.tmpl" }}<html></html>{{end}}`嵌套

  ```go
  func main() {
  	router := gin.Default()
  	router.LoadHTMLGlob("templates/**/*")
  	router.GET("/posts/index", func(c *gin.Context) {
  		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
  			"title": "Posts",
  		})
  	})
  	router.GET("/users/index", func(c *gin.Context) {
  		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
  			"title": "Users",
  		})
  	})
  	router.Run(":8080")
  }
  ```

  templates/posts/index.tmpl

  ```html
  # 对应router定义路径
  {{ define "posts/index.tmpl" }}
  <html>
      <h1>
  		{{ .title }}
  	</h1>
  	<p>Using posts/index.tmpl</p>
  </html>
  # 结束
  {{ end }}
  ```

  templates/users/index.tmpl

  ```html
  # 对应router定义路径
  {{ define "users/index.tmpl" }}
  <html>
      <h1>
  		{{ .title }}
  	</h1>
  	<p>Using users/index.tmpl</p>
  </html>
  # 结束
  {{ end }}
  ```

  访问：

  ```shell
  GET http://127.0.0.1:8080/posts/index
  GET http://127.0.0.1:8080/users/index
  ```

  

## 自定义模板渲染器

 你可以使用自定义的 html 模板渲染 

```go
package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func main() {
	router := gin.Default()
	html := template.Must(template.ParseFiles("templates/users.tmpl", "templates/posts.tmpl"))
	router.SetHTMLTemplate(html)
	//router.Run(":8080")

	router.GET("/posts/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
		})
	})
	router.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "Users",
		})
	})

	router.Run(":8080")
}
```

templates/users.tmpl

```html
{{ define "posts/index.tmpl" }}
<html>
    <h1>{{ .title }}</h1>
    <p>Using post.tmpl</p>
</html>
{{ end }}
```

templates/posts.tmpl

```html
{{ define "users/index.tmpl" }}
<html>
    <h1>{{ .title }}</h1>
    <p>Using users.tmpl</p>
</html>
{{ end }}
```

## 定义分隔符

你可以使用自定义分隔,必须写在模版渲染前

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	// 自定义分隔符
	router.Delims("{[{", "}]}")
	router.LoadHTMLGlob("templates/users.tmpl")

	router.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users.tmpl", gin.H{
			"title": "Posts",
		})
	})

	router.Run(":8080")
}
```

templates/users.html

```html
<html>
    <h1>{[{ .title }]}</h1>
    <p>Using users.tmpl</p>
</html>
```

## 自定义模板功能

* `LoadHTMLFiles() `

  ```go
  package main
  
  import (
  	"fmt"
  	"html/template"
  	"time"
  	"github.com/gin-gonic/gin"
  	"net/http"
  )
  
  func formatAsDate(t time.Time) string {
  	year, month, day := t.Date()
  	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
  }
  
  func main() {
  	router := gin.Default()
  	router.Delims("{[{", "}]}")
  	router.SetFuncMap(template.FuncMap{
  		"formatAsDate": formatAsDate,
  	})
  	router.LoadHTMLFiles("./templates/raw.tmpl")
  
  	router.GET("/raw", func(c *gin.Context) {
  		c.HTML(http.StatusOK, "raw.tmpl", map[string]interface{}{
  			"now": time.Date(2017, 07, 01, 0, 0, 0, 0, time.UTC),
  		})
  	})
  
  	router.Run(":8080")
  }
  ```

  raw.tmpl
  
  ```go
  Date: {[{.now | formatAsDate}]}
  ```
  
  

# 上传文件

**multipart forms 设置较低的内存限制 (默认是 32 MiB)**，并不是限制上传文件大小

自定义上传文件大小 1<< 20 为 1MB

```go
router.MaxMultipartMemory = 8 << 20  // 8 MiB
```

## 上传单个文件

```go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const uploaddir = "uploaddir"
func main() {
	router := gin.Default()
	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/upload", upload)
	router.Run(":8080")
}

//单文件上传
func upload(c *gin.Context) {
	//获取文件头
	file, err := c.FormFile("file-file2")
	if err != nil {
		c.String(http.StatusBadRequest, "请求失败")
		return
	}
	//获取文件名
	fileName := file.Filename
	fmt.Println("文件名：", fileName)
	//保存文件到服务器本地
	//SaveUploadedFile(文件头，保存路径)
	if err = c.SaveUploadedFile(file, uploaddir+"/"+fileName); err != nil {
		c.String(http.StatusBadRequest, "保存失败 Error:%s", err.Error())
		return
	}
	c.String(http.StatusOK, "上传文件成功")
}
```

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
<form method="post" action="http://127.0.0.1:8080/upload" enctype='multipart/form-data'>
    <input type="file" name="file-file2">
    <input type="submit" name="提交">
</form>
</body>
</html>
```

测试：

```shell
$ curl -X POST http://localhost:8080/upload \
  -F "file=@/Users/appleboy/test.zip" \
  -H "Content-Type: multipart/form-data"
  
 文件上传成功
```



## 上传多个文件

```go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const uploaddir = "uploaddir"
func main() {
	router := gin.Default()
	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/upload", upload)
	router.Run(":8080")
}

//单文件上传
func upload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "请求失败")
		return
	}
	files := form.File["file-file2[]"]

	for _,file := range files{
		log.Println(file.Filename)

		// 上传文件至指定目录
		_ = c.SaveUploadedFile(file,uploaddir+"/"+file.Filename)
	}
	c.String(http.StatusOK,fmt.Sprintf("%d files uploaded!", len(files)))
}
```

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>

<form method="post" action="http://127.0.0.1:8080/upload" enctype='multipart/form-data'>
    <input type="file" name="file-file2[]" multiple>
    <input type="submit" name="提交">
</form>

</body>
</html>
```

测试：

```shell
$ curl -X POST http://localhost:8080/upload \
  -F "upload[]=@/Users/appleboy/test1.zip" \
  -F "upload[]=@/Users/appleboy/test2.zip" \
  -H "Content-Type: multipart/form-data"
```



# form表单































