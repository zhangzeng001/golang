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

