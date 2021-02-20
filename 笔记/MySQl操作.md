# Go操作MySQL

## 连接

Go语言中的`database/sql`包提供了保证SQL或类SQL数据库的泛用接口，并不提供具体的数据库驱动。使用`database/sql`包时必须注入（至少）一个数据库驱动。

我们常用的数据库基本上都有完整的第三方实现。例如：[MySQL驱动](https://github.com/go-sql-driver/mysql)

### 下载依赖

`go get -u github.com/go-sql-driver/mysql`

### 使用MySQL驱动

```go
func Open(driverName, dataSourceName string) (*DB, error)
```

Open打开一个dirverName指定的数据库，dataSourceName指定数据源，一般至少包括数据库文件名和其它连接必要的信息。 

```go
package main

import(
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "root:111111@tcp(127.0.0.1:3306)/godb"
	db ,err := sql.Open("mysql",dsn)
	if err != nil{
		fmt.Println(err)
		return
	}
	defer db.Close()   // 注意这行代码要写在上面err判断的下面
}
```

### 初始化连接

Open函数可能只是验证其参数格式是否正确，实际上并不创建与数据库的连接。如果要检查数据源的名称是否真实有效，应该调用Ping方法。

返回的DB对象可以安全地被多个goroutine并发使用，并且维护其自己的空闲连接池。因此，Open函数应该仅被调用一次，很少需要关闭这个DB对象。

```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
func initDB() (err error) {
	// DSN:Data Source Name
	dsn := "root:111111@tcp(127.0.0.1:3306)/godb?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil{
		return err
	}
	return nil
}


func main() {
	err := initDB() // 调用输出化数据的函数
	if err != nil{
		fmt.Println(err)
		return
	}
}
```

其中`sql.DB`是表示连接的数据库对象（结构体实例），它保存了连接数据库相关的所有信息。它内部维护着一个具有零到多个底层连接的连接池，它可以安全地被多个goroutine同时使用。 

### SetMaxOpenConns

`func (db *DB) SetMaxOpenConns(n int)`

`SetMaxOpenConns`设置与数据库建立连接的最大数目。 如果n大于0且小于最大闲置连接数，会将最大闲置连接数减小到匹配最大开启连接数的限制。 如果n<=0，不会限制最大开启连接数，默认为0（无限制）。 

### SetMaxIdleConns

`func (db *DB) SetMaxIdleConns(n int)`

SetMaxIdleConns设置连接池中的最大闲置连接数。 如果n大于最大开启连接数，则新的最大闲置连接数会减小到匹配最大开启连接数的限制。 如果n<=0，不会保留闲置连接。 

```go
dsn := "root:111111@tcp(127.0.0.1:3306)/godb?charset=utf8mb4&parseTime=True"
// 不会校验账号密码是否正确
db, err = sql.Open("mysql", dsn)
// 设置与数据库建立连接的最大数目
db.SetMaxOpenConns(200)
// 设置连接池中的最大闲置连接数
db.SetMaxIdleConns(30)
```



## 增删改查 CURD

 不同的数据库中，SQL语句使用的占位符语法不尽相同。 

|   数据库   |  占位符语法  |
| :--------: | :----------: |
|   MySQL    |     `?`      |
| PostgreSQL | `$1`, `$2`等 |
|   SQLite   |  `?` 和`$1`  |
|   Oracle   |   `:name`    |

### 建库建表

我们先在MySQL中创建一个名为`sql_test`的数据库

```sql
CREATE DATABASE sql_test;
```

进入该数据库:

```sql
use sql_test;
```

执行以下命令创建一张用于测试的数据表：

```sql
CREATE TABLE `user` (
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(20) DEFAULT '',
    `age` INT(11) DEFAULT '0',
    PRIMARY KEY(`id`)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
```

### 查询 DQL

为了方便查询，我们事先定义好一个结构体来存储user表的数据。

```go
type user struct {
	id   int
	age  int
	name string
}
```

#### 单行查询

单行查询`db.QueryRow()`执行一次查询，并期望返回最多一行结果（即Row）。QueryRow总是返回非nil的值，直到返回值的Scan方法被调用时，才会返回被延迟的错误。（如：未找到结果） 

`func (db *DB) QueryRow(query string, args ...interface{}) *Row`

`columns, _ := row.Columns() 可以获取字段数组 func (rs *Rows) Columns() ([]string, error)`

```go
	//sqlStr := fmt.Sprintf(q)
	var u user
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	sqlStr := "select id, name, age from user where id=2"
	// 此处和以下一条有区别
	err := db.QueryRow(sqlStr).Scan(&u.id,&u.name,&u.age)
	if err != nil{
		fmt.Println("errrrrrrr")
		fmt.Println(err)
		return
	}
	fmt.Printf("id:%v name:%v age=%d",u.id,u.name,u.age)
```

```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
func initDB() (err error) {
	// DSN:Data Source Name
	dsn := "root:111111@tcp(127.0.0.1:3306)/godb?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	db, err = sql.Open("mysql", dsn)
	// 设置与数据库建立连接的最大数目
	//db.SetMaxOpenConns(200)
	// 设置连接池中的最大闲置连接数
	//db.SetMaxIdleConns(30)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil{
		return err
	}
	return nil
}

type user struct {
	id   int
	age  int
	name string
}

// 单行查询
func queryRowDemo(){
	//sqlStr := fmt.Sprintf(q)
	var u user
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	sqlStr := "select id, name, age from user where id=?"
	// args 对应 ？值
	err := db.QueryRow(sqlStr,2).Scan(&u.id,&u.name,&u.age)
	if err != nil{
		fmt.Println("errrrrrrr")
		fmt.Println(err)
		return
	}
	fmt.Printf("id:%v name:%v age=%d",u.id,u.name,u.age)
}

func main() {
	err := initDB() // 调用输出化数据的函数
	if err != nil{
		fmt.Println(err)
		return
	}

	queryRowDemo()
}
```

#### 多行查询

```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
func initDB() (err error) {
	// DSN:Data Source Name
	dsn := "root:111111@tcp(127.0.0.1:3306)/godb?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	db, err = sql.Open("mysql", dsn)
	// 设置与数据库建立连接的最大数目
	//db.SetMaxOpenConns(200)
	// 设置连接池中的最大闲置连接数
	//db.SetMaxIdleConns(30)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil{
		return err
	}
	return nil
}

type user struct {
	id   int
	age  int
	name string
}

// 单行查询
func queryRowDemo(){
	//sqlStr := fmt.Sprintf(q)
	var u user
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	sqlStr := "select id, name, age from user where id=?"
	// args 对应 ？值
	err := db.QueryRow(sqlStr,2).Scan(&u.id,&u.name,&u.age)
	if err != nil{
		fmt.Println("errrrrrrr")
		fmt.Println(err)
		return
	}
	fmt.Printf("id:%v name:%v age=%d",u.id,u.name,u.age)
}

// 多行查询 "select id ,name,age from user where id >2"
func queryMultiRowDemo(){
	sqlStr := "select id ,name,age from user where id >2"
	rows , err := db.Query(sqlStr)
	if err != nil{
		fmt.Println(err)
		return
	}
	// 关闭rows释放持有的数据库链接
	defer rows.Close()
	// 循环读取结果集中的数据
	for rows.Next(){
		var u user
		err := rows.Scan(&u.id,&u.name,&u.age)
		if err != nil{
			fmt.Println(err)
			return
		}
		fmt.Println(u.id,u.name,u.age)
	}
}
// 多行查询 "select * from user"
func queryMultiRowDemo2(){
	sqlStr := "select * from user"
	rows , err := db.Query(sqlStr)
	if err != nil{
		fmt.Println(err)
		return
	}
	// 关闭rows释放持有的数据库链接
	defer rows.Close()
	// 循环读取结果集中的数据
	for rows.Next(){
		var u user
		err := rows.Scan(&u.id,&u.name,&u.age)
		if err != nil{
			fmt.Println(err)
			return
		}
		fmt.Println(u.id,u.name,u.age)
	}
}
// 多行查询 函数传参
func queryMultiRowDemo3(sqlStr string,args ...interface{}){
	rows , err := db.Query(sqlStr,args...)
	if err != nil{
		fmt.Println(err)
		return
	}
	// 关闭rows释放持有的数据库链接
	defer rows.Close()
	// 循环读取结果集中的数据
	for rows.Next(){
		var u user
		err := rows.Scan(&u.id,&u.name,&u.age)
		if err != nil{
			fmt.Println(err)
			return
		}
		fmt.Println(u.id,u.name,u.age)
	}
}
func main() {
	err := initDB() // 调用输出化数据的函数
	if err != nil{
		fmt.Println(err)
		return
	}
    // 方式1
    qStr := `select * from user;`
    // 方式2
    qStr2 := `select * from user where id > ?`
    // 方式3
    qStr3 := `select * from user where id > ? and name = ?`
    id := 0
    name := "u4"
	queryMultiRowDemo3(qStr)
    fmt.Println("---------------------")
	queryMultiRowDemo3(qStr2,id)
	fmt.Println("---------------------")
	queryMultiRowDemo3(qStr3,id,name)
}

/*
1 zhang3 21
2 h2 21
3 u3 18
4 u4 19
---------------------
1 zhang3 21
2 h2 21
3 u3 18
4 u4 19
---------------------
4 u4 19
*/
```



### 增、删、改 DML

 插入、更新和删除操作都使用`Exec`方法。 

`func (db *DB) Exec(query string, args ...interface{}) (Result, error)`

* Result 有两个方法：

  `LastInsertId()`返回新插入数据的id

  `RowsAffected()`返回影响行数

```go
type Result interface {
	// LastInsertId returns the integer generated by the database
	// in response to a command. Typically this will be from an
	// "auto increment" column when inserting a new row. Not all
	// databases support this feature, and the syntax of such
	// statements varies.
	LastInsertId() (int64, error)

	// RowsAffected returns the number of rows affected by an
	// update, insert, or delete. Not every database or database
	// driver may support this.
	RowsAffected() (int64, error)
}
```

Exec执行一次命令（包括查询、删除、更新、插入等），返回的Result是对已执行的SQL命令的总结。参数args表示query中的占位参数。 

#### 插入数据

```go
// insertRow 插入数据
func insertRow(){
	// 要执行的sql语句
	//sqlstr := `insert into user(name ,age) values("wang5",22)`
	//sqlstr := `insert into user values(8,"yuan8",20)`
	sqlstr := "insert into user values(?,?,?)"
	// 执行sql,err返回执行错误
	res , err := db.Exec(sqlstr,9,"yang9",12)
	if err != nil{
		fmt.Println(err)
		return
	}
	// 得到新插入数据的id
	id , err := res.LastInsertId()
	if err != nil{
		fmt.Println(err)
		return
	}
	// 影响行数
	r, err := res.RowsAffected()
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Printf("Insert success, the id is %d.\n",id)
	fmt.Printf("%d rows in set.\n",r)
}
```

#### 删除数据

```go
func deleteRow(){
	sqlstr := `delete from user where id = ?`
	res,err := db.Exec(sqlstr,9)
	if err != nil{
		fmt.Println(err)
	}
	n , err := res.RowsAffected()
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Printf("Query OK, %d row affected.", n)
}
```

#### 修改数据

```go
func updateRow(){
	sqlstr := `update user set name=? where id = ?`
	res ,err := db.Exec(sqlstr,"yuan",8)
	if err != nil{
		fmt.Println(err)
		return
	}
	n, err := res.RowsAffected()
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Printf("Changed: %d",n)
}
```



## MySQL预处理

### 什么是预处理？

普通SQL语句执行过程：

1. 客户端对SQL语句进行占位符替换得到完整的SQL语句。
2. 客户端发送完整SQL语句到MySQL服务端
3. MySQL服务端执行完整的SQL语句并将结果返回给客户端。

预处理执行过程：

1. 把SQL语句分成两部分，命令部分与数据部分。
2. 先把命令部分发送给MySQL服务端，MySQL服务端进行SQL预处理。
3. 然后把数据部分发送给MySQL服务端，MySQL服务端对SQL语句进行占位符替换。
4. MySQL服务端执行完整的SQL语句并将结果返回给客户端。

### 为什么要预处理？

1. 优化MySQL服务器重复执行SQL的方法，可以提升服务器性能，提前让服务器编译，一次编译多次执行，节省后续编译的成本。
2. 避免SQL注入问题。

### Go实现MySQL预处理

`database/sql`中使用下面的`Prepare`方法来实现预处理操作。

```go
func (db *DB) Prepare(query string) (*Stmt, error)
```

`Prepare`方法会先将sql语句发送给MySQL服务端，返回一个准备好的状态用于之后的查询和命令。返回值可以同时执行多个查询和命令。

查询操作的预处理示例代码如下：

```go
// 预处理查询示例
func prepareQueryDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	defer rows.Close()
	// 循环读取结果集中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
	}
}
```

插入、更新和删除操作的预处理十分类似，这里以插入操作的预处理为例：

```go
// 预处理插入示例
func prepareInsertDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec("小王子", 18)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	_, err = stmt.Exec("沙河娜扎", 18)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	fmt.Println("insert success.")
}
```

### SQL注入问题

**我们任何时候都不应该自己拼接SQL语句！**

这里我们演示一个自行拼接SQL语句的示例，编写一个根据name字段查询user表的函数如下：

```go
// sql注入示例
func sqlInjectDemo(name string) {
	sqlStr := fmt.Sprintf("select id, name, age from user where name='%s'", name)
	fmt.Printf("SQL:%s\n", sqlStr)
	var u user
	err := db.QueryRow(sqlStr).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("exec failed, err:%v\n", err)
		return
	}
	fmt.Printf("user:%#v\n", u)
}
```

此时以下输入字符串都可以引发SQL注入问题：

```go
sqlInjectDemo("xxx' or 1=1#")
sqlInjectDemo("xxx' union select * from user #")
sqlInjectDemo("xxx' and (select count(*) from user) <10 #")
```



## Go实现MySQL事务

### 事务相关方法

Go语言中使用以下三个方法实现MySQL中的事务操作。 开始事务

```go
func (db *DB) Begin() (*Tx, error)
```

提交事务

```go
func (tx *Tx) Commit() error
```

回滚事务

```go
func (tx *Tx) Rollback() error
```

示例：

```go
// 事务操作示例
func transactionDemo() {
	tx, err := db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		fmt.Printf("begin trans failed, err:%v\n", err)
		return
	}
	sqlStr1 := "Update user set age=30 where id=?"
	ret1, err := tx.Exec(sqlStr1, 2)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql1 failed, err:%v\n", err)
		return
	}
	affRow1, err := ret1.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	sqlStr2 := "Update user set age=40 where id=?"
	ret2, err := tx.Exec(sqlStr2, 3)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql2 failed, err:%v\n", err)
		return
	}
	affRow2, err := ret2.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	fmt.Println(affRow1, affRow2)
	if affRow1 == 1 && affRow2 == 1 {
		fmt.Println("事务提交啦...")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		fmt.Println("事务回滚啦...")
	}

	fmt.Println("exec trans success!")
}
```



## 操作类封装

### 自定义类

* **README**

```go
# 初始化
​```go
// 拼接数据库连接信息字符串
dsn := "root:111111@tcp(127.0.0.1:3306)/godb?charset=utf8mb4&parseTime=True"
// 获取db对象
var db MyMysql
​```
# 查询
支持单行，多行查询，返回值未map的切片
`func (m *MyMysql)RowsQuery(sqlstr string, args ...interface{})([]map[string]string,error)`
sqlstr 要执行的sql语句传参
ages   sql语句中要替代的字符
​```go
// 支持类型
// 1  返回信息  [map[id:1] map[id:2]]
sqlstr := "select id from user where id < ? and id > ?"
res , err := db.RowsQuery(sqlstr,3,1)
// 2
sqlstr := "select id from user where id < ? "
res , err := db.RowsQuery(sqlstr,3)
// 3  [map[age:21 id:1 name:zhang3] map[age:21 id:2 name:li4]]
sqlstr := "select * from user where id < ? "
res , err := db.RowsQuery(sqlstr,3)
// 4  
sqlstr := "select * from user"
res , err := db.RowsQuery(sqlstr)
​```
```

* 自定义类

```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DbOperate interface {

}

// MyMysql 自定义结构体构建 DB对象
type MyMysql struct {
	Db *sql.DB
}

// Init 初始化数据库
func (m *MyMysql)Init(dsn string)(err error){
	// DSN:Data Source Name
	dsn = dsn
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	db, err := sql.Open("mysql", dsn)
	// 设置与数据库建立连接的最大数目
	//db.SetMaxOpenConns(200)
	// 设置连接池中的最大闲置连接数
	//db.SetMaxIdleConns(30)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil{
		return err
	}
	fmt.Println("连接数据库成功...")
	// 赋值给结构体，得到db对象
	m.Db = db
	return nil
}

// RowsQuery 查询，支持单行多行
func (m *MyMysql)RowsQuery(sqlstr string, args ...interface{})([]map[string]string,error){
	// 查询返回row对象
	rows,err := m.Db.Query(sqlstr, args...)
	if err != nil{
		fmt.Println(err)
		return nil,err
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()
	// 获取字段名信息 eg：[id name age]
	columns,err := rows.Columns()
	if err != nil{
		fmt.Println(err)
		return nil,err
	}
	// queryRes 创建返回数据容器
	var queryRes []map[string]string
	// scanArgs 接口类型适用于接收row.Scan(interface)返回
	scanArgs := make([]interface{},len(columns))
	// values 创建查询数据存储容器，构建返回map时interface类型不匹配map字段string类型，所以此处需要再定义values
	values := make([]string, len(columns))
	// scanArgs放入values元素的指针 eg：[0x11c34140 0x11c34148 0x11c34150]
	for i,_ := range scanArgs{
		scanArgs[i] = &values[i]
	}
	// 拿到查询到的具体数据 eg：scanargs [0x11c34140 0x11c34148 0x11c34150]；； values [3 zhang3 18]
	for rows.Next(){
		err = rows.Scan(scanArgs...)
		if err != nil{
			fmt.Println(err)
			return nil,err
		}

		r := make(map[string]string)
		// 构建返回数据容器
		for k,v := range values{
			column := columns[k]
			// 方式未初始化map引发panic
			//if queryRes[count] == nil{
			//	queryRes[count] = make(map[string]string)
			//}
			//queryRes[count][column] = v
			r[column] = v
		}
		queryRes = append(queryRes,r)
	}
	// 省去此处操作
	//dic := queryRes[0]
	//for i := 0 ;i < len(columns); i++{
	//	res := reflect.ValueOf(scanArgs[0]).String()
	//	dic[columns[i]] = res
	//}

	// 返回数据
	return queryRes,nil
}


```



* 自定义类封装操作

```
// 查询
dsn := "root:111111@tcp(127.0.0.1:3306)/godb?charset=utf8mb4&parseTime=True"
var db MyMysql
err := db.Init(dsn)
if err != nil{
panic(err)
}
sqlstr := "select * from user"
res , err := db.RowsQuery(sqlstr)
if err != nil{
panic(err)
}
defer db.Db.Close()
fmt.Println(res)
```



### github例子

[github](https://github.com/shangyou/gomysql/blob/master/mysql.go)

```go
package gomysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	Db *sql.DB
}

func (this *Mysql) Init(dns string) error {
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return err
	}

	this.Db = db
	return nil
}

func (this *Mysql) query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := this.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return rows, err
}

func (this *Mysql) Exec(query string, args ...interface{}) (sql.Result, error) {
	res, err := this.Db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *Mysql) Insert(query string, args ...interface{}) (int64, error) {
	res, err := this.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return id, err
}

func (this *Mysql) Update(query string, args ...interface{}) (int64, error) {
	res, err := this.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()

	return count, nil
}

func (this *Mysql) Delete(query string, args ...interface{}) (int64, error) {
	res, err := this.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()

	return count, err
}

func (this *Mysql) Fetchrow(query string, args ...interface{}) map[string]string {
	row, _ := this.query(query, args...)
	columns, _ := row.Columns()

	scanArgs := make([]interface{}, len(columns))
	values := make([]string, len(columns))

	for i, _ := range scanArgs {
		scanArgs[i] = &values[i]
	}

	for row.Next() {
		err := row.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}
	}
	rs := make(map[string]string)
	for i, _ := range values {
		rs[columns[i]] = values[i]
	}

	return rs
}

func (this *Mysql) Fetchrows(query string, args ...interface{}) map[int]map[string]string {
	row, _ := this.query(query, args...)
	columns, _ := row.Columns()
    // 接口类型适用于接收row.Scan(interface)返回
	scanArgs := make([]interface{}, len(columns))
    // 构建返回map时interface类型不匹配map字段类型
	values := make([]string, len(columns))
	
    // interface存values值的指针，scan接收interface能赋值给对应指针
	for i, _ := range scanArgs {
		scanArgs[i] = &values[i]
	}
    /*   省去类型转换操作步骤
    scanArgs := make([]interface{},len(columns))
    dic := queryRes[0]
    for i := 0 ;i < len(columns); i++{
		res := reflect.ValueOf(scanArgs[0]).String()
		dic[columns[i]] = res
	}
    */

	rs := make(map[int]map[string]string)
	i := 0
	for row.Next() {
		err := row.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		r := make(map[string]string)
		for k, v := range values {
			r[columns[k]] = v
		}

		rs[i] = r
		i++
	}

	return rs
}
```

