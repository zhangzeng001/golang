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

## 使用 [Govendor](https://github.com/kardianos/govendor) 工具创建项目(未使用)

1.`go get govendor`

```sh
$ go get github.com/kardianos/govendor
```

2.创建项目并且 `cd` 到项目目录中

```sh
$ mkdir -p $GOPATH/src/github.com/myusername/project && cd "$_"
```

3.使用 govendor 初始化项目，并且引入 gin

```sh
$ govendor init
$ govendor fetch github.com/gin-gonic/gin@v1.3
```

4.复制启动文件模板到项目目录中

```sh
$ curl https://raw.githubusercontent.com/gin-gonic/examples/master/basic/main.go > main.go
```

5.启动项目

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



# 示例

## AsciiJSON

```golang
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
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

  















