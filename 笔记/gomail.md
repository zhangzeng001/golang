# gomail

https://github.com/go-gomail/gomail

https://pkg.go.dev/gopkg.in/gomail.v2

## Download

```
go get gopkg.in/gomail.v2
```



## daemon

```go
package main

import (
	"gopkg.in/gomail.v2"
)

const (
	MailServerHost = "smtp.163.com"        // smtp地址
	MailServerPort = 465                   // smtp服务器端口
	MailServerUser = "1@163.com" // smtp发送用户
	MailServerpass = ""            // smtp用户密码
)

func main() {
	m := gomail.NewMessage()
	// 这种方式可以添加别名，即 nickname， 也可以直接用m.SetHeader("From", MAIL_USER)
	nickname := "发件人别名,使用中文会编码"
	m.SetHeader("From",nickname + "<" + MailServerUser + ">")
	// 多个收件人
	m.SetHeader("To", "343742221@qq.com","834148284@qq.com")
	// 抄送
	m.SetHeader("Cc", "834148284@qq.com")
	// 设置邮件主题
	m.SetHeader("Subject", "邮件名")
	// 设置邮件正文
	m.SetBody("text/html", "这是邮件内容")
	// 发送附件
	m.Attach("./node_exporter.tar.gz")
	d := gomail.NewDialer(MailServerHost, MailServerPort, MailServerUser, MailServerpass)
	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
```



## 邮件类



```go
package main

import (
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
)

const (
	MailServerHost     = "smtp.163.com"        // smtp地址
	MailServerPort     = 465                   // smtp服务器端口
	MailServerUser     = "1@163.com" // smtp发送用户
	MailServerpass     = ""            // smtp用户密码
	MailServerNickname = "zhang"               // 发件人别名，注意中文乱码
)

type Recv struct {
	To []string
	Cc []string
	Subject string
	body string
	Attach string
}

// SendMail 邮件发送
func (r Recv)SendMail()(msg string,err error){
	// 判断是否又收件人
	if len(r.To) == 0{
		err := errors.New("无收件人")
		return "",err
	}

	m := gomail.NewMessage()
	// 添加别名，也可以直接用m.SetHeader("From", MAIL_USER)
	nickname := MailServerNickname
	m.SetHeader("From",nickname + "<" + MailServerUser + ">")
	// 组织收件人切片
	m.SetHeader("To",r.To...)
	// 抄送
	if len(r.Cc) != 0{
		//fmt.Println(r.Cc)
		m.SetHeader("Cc",r.Cc...)
	}
	// 设置邮件主题
	m.SetHeader("Subject", r.Subject)
	// 设置邮件正文
	m.SetBody("text/html", r.body)
	// 发送附件
	if len(r.Attach) !=0{
		m.Attach(r.Attach)
	}

	d := gomail.NewDialer(MailServerHost, MailServerPort, MailServerUser, MailServerpass)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return "", err
	}
	return "success",nil
}

func main() {
	var r = Recv{
		To: []string{"343742221@qq.com","834148284@qq.com"},
		Subject: "测试邮件主题",
		body: "这是邮件内容",
	}
	res ,err := r.SendMail()
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(res)
}
```

![1613870033679](gomail.assets\1613870033679.png)

### 遇到问题

传入空字符串会报错，但能发送成功

因为r.To和...之间加了空格`m.SetHeader("To",r.To...)`

```go
gomail: could not send email 1: gomail: invalid address "": mail: no address
```



切片类型传参到函数为"a","b","c"

```go
package main

import "fmt"

// 函数接收，value传入会自动变成[]切片
func test(v string,value ...string)  {
	for k,v1 := range value{
		fmt.Println(k)
		fmt.Println(v1)
	}
}

func main() {
	s := []string{"343742221@qq.com","834148284@qq.com"}
    // 传入参数是直接... 如果切片和...之间有空格，会传入空字符串""
	test("a",s...)
}
```



# gomail发送邮件之解决附件名乱码的问题

https://www.jianshu.com/p/ab63ee725888

在使用Go语言开发时，我们会遇到发送邮件的需求，在Go语言标准包中，也提供了邮件发送客户端`smtp`的封装。不过，该标准包只提供了基础的邮件发送过程，对于一些复杂的定义还需要自己去封装，封装过程就需要依据邮件协议[RFC2822](https://tools.ietf.org/html/rfc2822)了。还好，github上有人专门为我们封装好了这个包：https://github.com/go-gomail/gomail。这个包封装了发送附件、图片、HTML内容模板、SSL和TLS等的支持，可以满足我们的大部分应用场景。下面，我就对`gomail`实现发送邮件做一下简单介绍。

### 1. 需要先安装`gomail`包



```csharp
$ go get -v gopkg.in/gomail.v2
```

### 2. 导入`gomail`包



```cpp
$ import "gopkg.in/gomail.v2"
```

### 3. 需要创建一个`Message`实例，`Message`提供了整个邮件协议内容的构建，默认实例采用`UTF-8`字符集和`Quoted-printable`编码。

> 对于`Quoted-printable`编码的定义，维基百科上是这样说的：Quoted-printable是使用可打印的ASCII字符（如字母、数字与“=”）表示各种编码格式下的字符，以便能在7-bit数据通路上传输8-bit数据, 或者更一般地说在非8-bit clean媒体上正确处理数。



```go
m := gomail.NewMessage()
```

### 4. 构造邮件内容，包括：发件人信息、收件人、主题、内容，更多内容设定可参考协议：[RFC2822](https://tools.ietf.org/html/rfc2822) 



```cpp
// 发件人信息
m.SetHeader("From", m.FormatAddress("user@example.com", "张三"))
// 收件人
m.SetHeader("To", "user@qq.com")
// 主题
m.SetHeader("Subject", "邮件标题")
// 内容
m.SetBody("text/html", "系统邮件请勿回复")
```

> 特殊说明，构造`From(发件人信息)`时需要使用`m.FormatAddress`方法，因为发件人指定中文名或特殊字符时，需要进行编码

### 5. 构造附件信息，同时对附件进行重命名

> 比如，我有一个临时文件：`/tmp/foo.txt`，我需要将这个文件以邮件附件的方式进行发送，同时指定附件名为：`附件.txt`



```go
name := "附件.txt"
m.Attach("/tmp/foo.txt",
    gomail.Rename(name),
)
```

### 6. 创建`SMTP`客户端，连接到远程的邮件服务器，需要指定服务器地址、端口号、用户名、密码，如果端口号为`465`的话，自动开启SSL，这个时候需要指定`TLSConfig` 

> 这里的用户名和密码指的是能够登录该邮箱的邮箱地址和密码



```go
d := gomail.NewDialer("smtp.example.com", 465, "user@example.com", "password")
d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
```

### 7. 执行邮件发送



```go
err := d.DialAndSend(m)
if err != nil {
    // 处理错误
}
```

至此，邮件已经发送成功了，整个邮件的内容为（其中，附件内容为`foo.bar`）：



```dart
Mime-Version: 1.0
Date: Sat, 10 Nov 2018 21:40:13 +0800
From: =?UTF-8?q?=E5=BC=A0=E4=B8=89?= <user@example.com>
To: user@qq.com
Subject: =?UTF-8?q?=E9=82=AE=E4=BB=B6=E6=A0=87=E9=A2=98?=
Content-Type: multipart/mixed;
 boundary=92142f9522a20d2f4feffce957bf68b46ad1a620e68fbecbf35669266e9a

--92142f9522a20d2f4feffce957bf68b46ad1a620e68fbecbf35669266e9a
Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=UTF-8

=E7=B3=BB=E7=BB=9F=E9=82=AE=E4=BB=B6=E8=AF=B7=E5=8B=BF=E5=9B=9E=E5=A4=8D
--92142f9522a20d2f4feffce957bf68b46ad1a620e68fbecbf35669266e9a
Content-Disposition: attachment; filename="附件.txt"
Content-Transfer-Encoding: base64
Content-Type: text/plain; charset=utf-8; name="附件.txt"

Zm9vLmJhcgo=
--92142f9522a20d2f4feffce957bf68b46ad1a620e68fbecbf35669266e9a--
```

> 打印邮件内容，可以将`Message`写入到一个缓冲区中，代码如下：



```dart
buf := new(bytes.Buffer)
m.WriteTo(buf)
fmt.Println(buf.String())
```

### 解决`gomail v2.0.0`版本下中文附件名乱码的问题

> 在不同的邮件服务器中，对于中文附件名的编码存在不同的规范，我们可以尝试一下，将上面的邮件附件发送到QQ邮箱，附件名显示正常，发送到126的邮箱就是乱码(这是我测试的结果)。对此，我们可以通过给附件名进行编码的方式来解决这个问题。



```go
    name := "附件.txt"
    m.Attach("/tmp/foo.txt",
        gomail.Rename(name),
        gomail.SetHeader(map[string][]string{
            "Content-Disposition": []string{
                fmt.Sprintf(`attachment; filename="%s"`, mime.QEncoding.Encode("UTF-8", name)),
            },
        }),
    )
    
```

### 将邮件内容更改为`Base64`编码



```go
    m := gomail.NewMessage(
        gomail.SetEncoding(gomail.Base64),
    )

// ...


    name := "附件.txt"
    m.Attach("/tmp/foo.txt",
        gomail.Rename(name),
        gomail.SetHeader(map[string][]string{
            "Content-Disposition": []string{
                fmt.Sprintf(`attachment; filename="%s"`, mime.BEncoding.Encode("UTF-8", name)),
            },
        }),
    )
```

> 使用`Base64`编码后的邮件内容为：



```dart
Mime-Version: 1.0
Date: Sat, 10 Nov 2018 21:53:22 +0800
From: =?UTF-8?b?5byg5LiJ?= <user@example.com>
To: user@qq.com
Subject: =?UTF-8?b?6YKu5Lu25qCH6aKY?=
Content-Type: multipart/mixed;
 boundary=42839966777f27ebe3861a73eabbf615036ea57b3222968e21519c64cdd5

--42839966777f27ebe3861a73eabbf615036ea57b3222968e21519c64cdd5
Content-Transfer-Encoding: base64
Content-Type: text/html; charset=UTF-8

57O757uf6YKu5Lu26K+35Yu/5Zue5aSN
--42839966777f27ebe3861a73eabbf615036ea57b3222968e21519c64cdd5
Content-Disposition: attachment; filename="=?UTF-8?b?6ZmE5Lu2LnR4dA==?="
Content-Transfer-Encoding: base64
Content-Type: text/plain; charset=utf-8; name="附件.txt"

Zm9vLmJhcgo=
--42839966777f27ebe3861a73eabbf615036ea57b3222968e21519c64cdd5--
```