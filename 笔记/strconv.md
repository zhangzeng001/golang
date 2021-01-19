# strconv包

`strconv`包实现了基本数据类型与其字符串表示的转换，主要有以下常用函数： `Atoi()`、`Itoa()`、parse系列、format系列、append系列。  



## string与int类型转换

 C语言中没有string类型而是用字符数组(array)表示字符串，所以`Itoa`对很多C系的程序员很好理解。 

### Atoi()

 `Atoi()`函数用于将字符串类型的整数转换为int类型，函数签名如下。 

`func Atoi(s string) (i int, err error)`

 如果传入的字符串参数无法转换为int类型，就会返回错误。 

```go
// 字符串转int
s1 := "100"
i1 ,err := strconv.Atoi(s1)
if err != nil{
    fmt.Println(err)
}
fmt.Printf("%T %v",i1,i1)  // int 100
```

### Itoa

`Itoa()`函数用于将int类型数据转换为对应的字符串表示，具体的函数签名如下。 

`func Itoa(i int) string`

```go
// int转字符串
i1 := 100
s1 := strconv.Itoa(i1)
fmt.Printf("%T %v",s1,s1)
```



## Parse系列函数

Parse类函数用于转换字符串为给定类型的值：ParseBool()、ParseFloat()、ParseInt()、ParseUint()。

### ParseBool()

 返回字符串表示的bool值。它接受1、0、t、f、T、F、true、false、True、False、TRUE、FALSE；否则返回错误。 

`func ParseBool(str string) (value bool, err error)`

```go
s1 := "true"
s2 := "t"
s3 := "F"
fmt.Println(strconv.ParseBool(s1))  // true <nil>
fmt.Println(strconv.ParseBool(s2))  // true <nil>
fmt.Println(strconv.ParseBool(s3))  // false <nil>
```



### ParseInt()

`func ParseInt(s string, base int, bitSize int) (i int64, err error)`

返回字符串表示的整数值，接受正负号。

base指定进制（2到36），如果base为0，则会从字符串前置判断，”0x”是16进制，”0”是8进制，否则是10进制；

bitSize指定结果必须能无溢出赋值的整数类型，0、8、16、32、64 分别代表 int、int8、int16、int32、int64；

返回的err是*NumErr类型的，如果语法有误，err.Error = ErrSyntax；如果结果超出类型范围err.Error = ErrRange。

```go
s1 := "100"
s2 , err := strconv.ParseInt(s1,10,32)
if err != nil{
    fmt.Println(err)
    return
}
fmt.Printf("%T %v",s2,s2) // int64 100
```



### ParseUnit()

`func ParseUint(s string, base int, bitSize int) (n uint64, err error)`

 `ParseUint`类似`ParseInt`但不接受正负号，用于无符号整型。 

```go
s1 := "100"
s2 , err := strconv.ParseUint(s1,10,32)
if err != nil{
    fmt.Println(err)
    return
}
fmt.Printf("%T %v",s2,s2) // int64 100
```



### ParseFloat()

`func ParseFloat(s string, bitSize int) (f float64, err error)`

解析一个表示浮点数的字符串并返回其值。

如果s合乎语法规则，函数会返回最为接近s表示值的一个浮点数（使用IEEE754规范舍入）。

bitSize指定了期望的接收类型，32是float32（返回值可以不改变精确值的赋值给float32），64是float64；

返回值err是*NumErr类型的，语法有误的，err.Error=ErrSyntax；结果超出表示范围的，返回值f为±Inf，err.Error= ErrRange。

```go
s1 := "2.2"
s2 , err := strconv.ParseFloat(s1,8)
if err != nil{
    fmt.Println(err)
    return
}
fmt.Printf("%T %v",s2,s2) // float64 2.2
```



## Format系列函数

### FormatBool()

`func FormatBool(b bool) string`

根据b的值返回”true”或”false”。 

```go
s1 := false
s2 := strconv.FormatBool(s1)
fmt.Printf("%T %v",s2,s2)  // string false
```



### FormatInt()

使用即可 `strconv.Itoa()`

`func FormatInt(i int64, base int) string`



### FormatUint()

使用即可 `strconv.Itoa()`

`func FormatUint(i uint64, base int) string`



### FormatFloat()

`func FormatFloat(f float64, fmt byte, prec, bitSize int) string`

函数将浮点数表示为字符串并返回。

bitSize表示f的来源类型（32：float32、64：float64），会据此进行舍入。

fmt表示格式：’f’（-ddd.dddd）、’b’（-ddddp±ddd，指数为二进制）、’e’（-d.dddde±dd，十进制指数）、’E’（-d.ddddE±dd，十进制指数）、’g’（指数很大时用’e’格式，否则’f’格式）、’G’（指数很大时用’E’格式，否则’f’格式）。

prec控制精度（排除指数部分）：对’f’、’e’、’E’，它表示小数点后的数字个数；对’g’、’G’，它控制总的数字个数。如果prec 为-1，则代表使用最少数量的、但又必需的数字来表示f。

```go
s1 := 3.1431
fmt.Println(strconv.FormatFloat(s1,'E',-1,64))
```

