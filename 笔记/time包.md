# time包

## 时间类型 time.Time

`time.Time`类型表示时间，通过`time.Now()`函数获取当前的时间对象，通过这个时间对象再去获取时、分、秒信息



```go
	// 获取当前时间对象
	now := time.Now()
	fmt.Printf("type:%T ; current time:%v",now,now)
	// time.Time
	// 2021-01-13 09:34:44.2185025 +0800 CST m=+0.014991101

	// 格式化输出时间
	strfmt:=now.Format("2006-01-02 15:04:05")

	// 通过时间对象获取具体时间
	year := now.Year()       // 年 2021
	mouth := now.Month()     // 月 January
	day := now.Day()         // 日 13
	hour := now.Hour()       // 时 9
	minute := now.Minute()   // 分 38
	second := now.Second()   // 秒 52

	fmt.Println(year,mouth,day,hour,minute,second)
```



## 时间戳

 时间戳是自1970年1月1日（08:00:00GMT）至当前时间的总毫秒数。它也被称为Unix时间戳（UnixTimestamp）。 

基于时间对象获取时间戳

```go
 	// 获取当前时间对象
	now := time.Now()

	timestamp1 := now.Unix()      // 时间戳     1610504202
	timestamp2 := now.UnixNano()  // 纳秒时间戳  1610504202961160300

	fmt.Println(timestamp1)
	fmt.Println(timestamp2)

	/*************************************/
	//时间转换,time.Unix()得到时间对象
	timeObj := time.Unix(1610504202,0)
	fmt.Println(timeObj)        // 2021-01-13 10:16:42 +0800 CST
	// 通过该对象再转换为具体时间
	year := timeObj.Year()      // 2021
	mouth := timeObj.Month()    // January
	day := timeObj.Day()        // 13
	hour := timeObj.Hour()      // 10
	minute := timeObj.Minute()  // 16
	second := timeObj.Second()  // 42

	fmt.Println(year,mouth,day,hour,minute,second)
```



## 时间间隔

`time.Duration`是`time`包定义的一个类型，它代表两个时间点之间经过的时间，以纳秒为单位。`time.Duration`表示一段时间间隔，可表示的最长时间段大约290年。 

 time包中定义的时间间隔类型的常量如下： 

```go
const (
    Nanosecond  Duration = 1
    Microsecond          = 1000 * Nanosecond
    Millisecond          = 1000 * Microsecond
    Second               = 1000 * Millisecond
    Minute               = 60 * Second
    Hour                 = 60 * Minute
)
```

 例如：`time.Duration`表示1纳秒，`time.Second`表示1秒。 



## 格式化时间

时间类型有一个自带的方法`Format`进行格式化，需要注意的是Go语言中格式化时间模板不是常见的`Y-m-d H:M:S`而是使用Go的诞生时间2006年1月2号15点04分（记忆口诀为2006 1 2 3 4）。也许这就是技术人员的浪漫吧。

补充：如果想格式化为12小时方式，需指定`PM`。

```go
	now := time.Now()
	// 24小时制
	fmt.Println(now.Format("2006-01-02 15:04:05.000 Mon Jan"))
	// 12小时制
	fmt.Println(now.Format("2006-01-02 03:04:05.000 PM Mon Jan"))
	fmt.Println(now.Format("2006/01/02 15:04"))
	fmt.Println(now.Format("15:04 2006/01/02"))
	fmt.Println(now.Format("2006/01/02"))
```

系统预定义格式：

```go
const (
	ANSIC       = "Mon Jan _2 15:04:05 2006"
	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = "02 Jan 06 15:04 MST"
	RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = "3:04PM"
	// Handy time stamps.
	Stamp      = "Jan _2 15:04:05"
	StampMilli = "Jan _2 15:04:05.000"
	StampMicro = "Jan _2 15:04:05.000000"
	StampNano  = "Jan _2 15:04:05.000000000"
)
```





## 时区

### 时区类型

```
loc值：Asia/Shanghai 类型：*time.Location
```

```go
	// 默认UTC
	loc1, _ := time.LoadLocation("")                      // UTC
	// 一般为CST
	loc2, _ := time.LoadLocation("Local")                 // Local  
	// 美国洛杉矶PDT
	loc3, _ := time.LoadLocation("America/Los_Angeles")   // America/Los_Angeles
	// CST
	loc4, _:= time.LoadLocation("Asia/Chongqing")         // Asia/Chongqing

	fmt.Println(loc1,loc2,loc3,loc4)
```



### 字符串转时间对象,非时区转换

`time.parse()`默认为UTC类型

```go
const TIME_LAYOUT = "2006-01-02 15:04:05"

func main() {
	res, err := time.Parse(TIME_LAYOUT, "2008-01-02 22:22:22")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%T\n",res)  // time.Time
	fmt.Println(res)        // 2008-01-02 22:22:22 +0000 UTC
}
```



转换前   `2021-01-13 11:27:41`

转换后   `2021-01-13 11:27:41 +0800 CST`



```go
now := time.Now()
fmt.Println(now)
// 加载时区
loc, err := time.LoadLocation("Asia/Shanghai")
if err != nil {
	fmt.Println(err)
	return
}
// 按照指定时区和指定格式解析字符串时间
timeObj, err := time.ParseInLocation("2006/01/02 15:04:05", "2019/08/04 14:15:20", loc)
if err != nil {
	fmt.Println(err)
	return
}
fmt.Println(timeObj)
fmt.Println(timeObj.Sub(now))
```



```go
package main

import (
	"fmt"
	"time"
)

const TIME_LAYOUT = "2006-01-02 15:04:05"

func parseWithLocation(name string, timeStr string) (time.Time, error) {
	locationName := name
	if l, err := time.LoadLocation(locationName); err != nil {
		return time.Time{}, err
	} else {
		//转成带时区的时间
		lt, _ := time.ParseInLocation(TIME_LAYOUT, timeStr, l)
		//直接转成相对时间
		//fmt.Println(time.Now().In(l).Format(TIME_LAYOUT))
		return lt, nil
	}
}

func main() {
	//testTime()
	str:=time.Now().Format(TIME_LAYOUT)
	//指定时区
	t1,_:=parseWithLocation("Asia/Shanghai", str)
	fmt.Println(t1)
}
```



### 时区转换

字符串时间  --> 先转换为该时区对应的格式 --> 再计算时间偏移量

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// 假设有一个字符串时间  2021-01-13 13:42:50.9754395 +0000 CST m=+0.013992401
	//timeStr := "2021-01-13 05:42:05.9754395 +0000 UTC"
	timeStr := "2021-01-13 05:42:05.9754395"

	// 格式化时间格式
	timeLayout := "2006-01-02 15:04:05"

	// 1.得到一个时区类型*time.Location
	loc ,err := time.LoadLocation("Asia/Shanghai") // loc值：Asia/Shanghai 类型：*time.Location
	if err != nil{
		fmt.Println(err)
		return
	}
	// 2.直接转换为该时区的时间对象
	theTime, err := time.ParseInLocation(timeLayout, timeStr, loc) //使用模板在对应时区转化为time.time类型
	if err != nil{
		fmt.Println(err)
		return
	}

	// 3.加减对应时区的时间偏移量 转为实际时区的具体时间
	utcstr := theTime.Add(+8 * time.Hour)
	fmt.Println(utcstr)   // 2021-01-13 13:42:05.9754395 +0800 CST
}
```



## 定时器

使用`time.Tick(时间间隔)`来设置定时器，定时器的本质上是一个通道（channel）。

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(1*time.Second)
	i:=0
	for{
		i++
		fmt.Println(<-ticker.C) // 输出时间，可以不打印直接写<-ticker.C
		if i==10{
			break
		}
	}
}
```



```go
func tickDemo() {
	ticker := time.Tick(time.Second) //定义一个1秒间隔的定时器
	for i := range ticker {
		fmt.Println(i)//每秒都会执行的任务
	}
}
```



# 常用时间操作

* 字符串转时间对象

  



* 获取当前格式化时间





## 时间加减 Add

 我们在日常的编码过程中可能会遇到要求时间+时间间隔的需求，Go语言的时间对象有提供Add方法如下： 

`func (t Time) Add(d Duration) Time`

```go
	// 当前时间加一个小时
	now := time.Now()
	a1 := now.Add(1 * time.Hour) // 2021-01-13 11:39:10.1425563 +0800 CST m=+3600.016990801
	fmt.Println(a1)
	// 当前时间减一个小时
	a2 := now.Add(-1 * time.Hour) // 2021-01-13 11:39:10.1425563 +0800 CST m=+3600.016990801
	fmt.Println(a2)
	// 某个时间加2天的时间
	a3 := now.Add(2 * 24 * time.Hour) // 2021-01-13 11:39:10.1425563 +0800 CST m=+3600.016990801
	fmt.Println(a3)
	// 某个时间减一天的时间
	a4 := now.Add(-2 * 24 * time.Hour) // 2021-01-13 11:39:10.1425563 +0800 CST m=+3600.016990801
	fmt.Println(a4)
	// 当前时间加一个月；注意一个月区分大小月
	a5 := now.Add(30 * 24  * time.Hour) // 2021-01-13 11:39:10.1425563 +0800 CST m=+3600.016990801
	fmt.Println(a5)
	// 当前是按减一个月
	a6 := now.Add(-30 * 24  * time.Hour) // 2021-01-13 11:39:10.1425563 +0800 CST m=+3600.016990801
	fmt.Println(a6)
```



## 时间差值 Sub

时间间隔

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	//fmt.Println(now)
	res := now.Sub(now.Add(1*time.Hour))
	fmt.Println(res)            // -1h0m0s
	fmt.Println(res.String())   // -1h0m0s
	fmt.Println(res.Seconds())  // -3600
	fmt.Println(res.Hours())    // -1
}
```



## 时间相同Equal（注意区分时区）

 判断两个时间是否相同，会考虑时区的影响，因此不同时区标准的时间也可以正确比较。本方法和用t==u不同，这种方法还会比较地点和时区信息。 

```go
	now := time.Now()
	res := now.Equal(time.Now())  // bool
	if res{
		fmt.Println("相同")
	}else {
		fmt.Println("不同")
	}
```



## Before

```go
func (t Time) Before(u Time) bool
```

如果t代表的时间点在u之前，返回真；否则返回假。

```go
	now := time.Now()
	res := now.Before(time.Now().Add(-24 * time.Hour))  // bool
	if res{
		fmt.Println("小于")
	}else {
		fmt.Println("大于")
	}
```



## After

```go
func (t Time) After(u Time) bool
```

如果t代表的时间点在u之后，返回真；否则返回假。



## 时间比较

```go
package main

import (
	"fmt"
	"time"
)

type mytimer struct {
	logtime time.Time
}

func main() {
	now := time.Now()
	var time1 = mytimer{logtime: now}
	for{
		// 判断当前时间是否在定义时间time1之后的6秒
		res := time1.logtime.Add(6 * time.Second).Before(time.Now())
		fmt.Println(res)
		fmt.Println("before--------------",time.Now())
		fmt.Println("struct--------------",time1.logtime)
		if res {
			time1.logtime = time.Now()
		}
		time.Sleep(1*time.Second)
	}
}
```



## 时间类型转换

```go
package main

import (
	"fmt"
	"time"
)

func main()  {

	Str2Time:=Str2Time("2017-09-12 12:03:40")
	fmt.Println(Str2Time)


	Str2Stamp:=Str2Stamp("2017-09-12 12:03:40")
	fmt.Println(Str2Stamp)

	Time2Str:=Time2Str()
	fmt.Println(Time2Str)

	GetStamp:=Time2Stamp()
	fmt.Println(GetStamp)

	Stamp2Str:=Stamp2Str(1505189020000)
	fmt.Println(Stamp2Str)

	Stamp2Time:=Stamp2Time(1505188820000)
	fmt.Println(Stamp2Time)
}

/**字符串->时间对象*/
func Str2Time(formatTimeStr string) time.Time{
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, formatTimeStr, loc) //使用模板在对应时区转化为time.time类型

	return theTime

}
/**字符串->时间戳*/
func Str2Stamp(formatTimeStr string) int64 {
	timeStruct:=Str2Time(formatTimeStr)
	millisecond:=timeStruct.UnixNano()/1e6
	return  millisecond
}

/**时间对象->字符串*/
func Time2Str() string {
	const shortForm = "2006-01-01 15:04:05"
	t := time.Now()
	temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	str := temp.Format(shortForm)
	return str
}


/*时间对象->时间戳*/
func Time2Stamp()int64{
	t:=time.Now()
	millisecond:=t.UnixNano()/1e6
	return  millisecond
}
/*时间戳->字符串*/
func Stamp2Str(stamp int64) string{
	timeLayout := "2006-01-02 15:04:05"
	str:=time.Unix(stamp/1000,0).Format(timeLayout)
	return str
}
/*时间戳->时间对象*/
func Stamp2Time(stamp int64)time.Time{
	stampStr:=Stamp2Str(stamp)
	timer:=Str2Time(stampStr)
	return timer
}
```

