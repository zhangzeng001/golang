package main

import (
	"awesomeProject/mylogger"
	"fmt"
	"time"
)

var lv = "info"
var console = false

// 接口方式
var mylog mylogger.LogInterface

func main() {
	var err error
	// 接口方式
	mylog ,err = mylogger.NewlogObj(lv,console,"")
	if err != nil{
		fmt.Println(err)
		return
	}
	for  {
		mylog.Debug("这是一条debug日志")
		mylog.Trace("这是一条Trace日志")
		mylog.Info("这是一条Info日志")
		mylog.Warning("这是一条Warning日志")
		mylog.Error("这是一条Error日志")
		mylog.Fatal("这是一条Fatal日志")

		time.Sleep(1*time.Second)
	}

}