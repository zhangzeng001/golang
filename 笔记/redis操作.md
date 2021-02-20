# Redis 操作

[文档](https://pkg.go.dev/github.com/go-redis/redis)

[Manager command](https://pkg.go.dev/github.com/go-redis/redis#Client.Command)

# go-redis库

## 安装

区别于另一个比较常用的Go语言redis client库：[redigo](https://github.com/gomodule/redigo)，这里采用https://github.com/go-redis/redis连接Redis数据库并进行操作，因为`go-redis`支持连接哨兵及集群模式的Redis。

使用以下命令下载并安装:

```bash
go get -u github.com/go-redis/redis
```

## 连接

### 普通连接

```go
// 声明一个全局的rdb变量
var rdb *redis.Client

// 初始化连接
func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
```

















