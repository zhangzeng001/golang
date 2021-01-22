# 正则表达式说明

Go正则表达式说明：

1. 大写英文字母的正则表达式，处理可以写成`[A-Z]`还可以写成`[\x41-\x5A]`因为在ASCII码字典中A-Z被排在了65-90号（也就是ASCII码的第66位到第91位），换算成十六禁止就是`\x41-\x5A`

2. `[0-9]`可以写成`[\x30-\x39]`

3. `[a-z]`可以写成`[\x61-\x7A]`

4. 中文正则表达式：`[\u4E00-\u9FA5]`

   因为中文在Unicode编码字典排在`u4E00-u9FA5`。换算成十进制也就是19968号到40869号是中文，一共2092个中文字体被收录到Unicode编码字典中



# 常用方法

## Match

检查正则表达式与**字节数组**是否匹配。更复杂的查询建议使用`regexp.Compile()`和更完整的regexp接口

`func Match(pattern string, b []byte) (matched bool, err error)`

```go
v1 := "123456789"
flag,_ := regexp.Match("^\\d{6,9}$",[]byte(v1))
flag2,_ := regexp.Match("^\\d{6,8}$",[]byte(v1))
fmt.Println(flag) // true
fmt.Println(flag2) // false
```

### MatchString

检查正则表达式与字符串是否匹配

`func MatchString(pattern string, s string) (matched bool, err error)`

```go
v1 := "123456789"
v2 := "abcdeFGh_"
flag,_ := regexp.MatchString("^\\d{6,15}$",v1)
flag2,_ := regexp.MatchString("^\\w",v2)
fmt.Println(flag)   // true
fmt.Println(flag2)  // true
```

## Compile

使用正则表达式对象做匹配

`func Compile(expr string) (*Regexp, error)`

```go
v1 := "123456789"
// 返回一个正则表达式对象
flag,err := regexp.Compile("^\\d{5,9}$")
if err != nil{
    fmt.Println(err)
    return
}
fmt.Println(flag.MatchString(v1)) // true
```

## MustCompile

和Compile用法相同，不同的是，正则表达式不能解析 不返回err，错误直接painc

`func MustCompile(str string) *Regexp`

```go
v1 := "123456789"
v2 := "1234sdfsdf"
// 返回一个正则表达式对象
flag := regexp.MustCompile("^\\d{3,9}$")
fmt.Println(flag.MatchString(v1))  // true
fmt.Println(flag.MatchString(v2))  // false
///////////////////////////////////////////////////////
v1 := "Fe匹配中文字符abc"
v2 := "匹配中文字符"
// 返回一个正则表达式对象
flag := regexp.MustCompile("^[\u4E00-\u9FA5]+$")
fmt.Println(flag.MatchString(v2))  // true
fmt.Println(flag.MatchString(v1))  // false

// 匹配数组
fmt.Println(flag.Match([]byte(v2)))  // true
```

## ReplaceAll和ReplaceAll

`func (re *Regexp) ReplaceAllString(src, repl string) string`

`ReplaceAllString` 针对字符串

`func (re *Regexp) ReplaceAll(src, repl []byte) []byte`

`ReplaceAll` 针对数组

```go
v1 := "Fe匹配中文字符abc"
// 返回一个正则表达式对象
regObj := regexp.MustCompile("[a-z]+")
ret := regObj.ReplaceAll([]byte(v1),[]byte("X"))
ret2 := regObj.ReplaceAllString(v1,"W")
fmt.Println(string(ret))  // FX匹配中文字符X
fmt.Println(string(ret2))  // FW匹配中文字符W
```

## Split

`func (re *Regexp) Split(s string, n int) []string`

将字符串按照正则表达式分割成子字符串组成的切片，如果切片长度超过指定参数n，则不再切分

```go
v1 := "192.168.10.11"
v2 := "abcd"
// 返回一个正则表达式对象
regObj := regexp.MustCompile("\\.+")
regObj2 := regexp.MustCompile("")
res := regObj.Split(v1,10)
res2 := regObj2.Split(v2,10)
fmt.Printf("%T   %v\n",res,res)  // []string   [192 168 10 11]
fmt.Printf("%T   %v",res2,res2)  // []string   [a b c d]
```

