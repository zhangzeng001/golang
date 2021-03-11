# Golang如何打包在Linux上部署

一、本地编译Golang程序

cmd控制台到main.go文件目录下，执行：

```shell
set GOARCH=amd64

set GOOS=linux

go build main.go

此时会生成一个没有后缀的二进制文件  main
```



二、上传Golang二进制文件到Linux服务器

将该文件放入linux系统某个文件夹下

赋予权限

```shell
chmod 777 main
最后执行 ./main 就行了。
```

如果想让项目在后台执行：执行 nohup ./main & ，这样就可以程序在后台运行了。 



-

-

-

-