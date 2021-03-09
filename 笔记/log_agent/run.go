package main

import (
	"awesomeProject/log_agent/elasPkg"
	"encoding/json"
	"fmt"
	"github.com/hpcloud/tail"
	"go.uber.org/zap"
	"github.com/panjf2000/ants/v2"
	"strings"
	"time"
)

func threadPool(f func(data interface{}),capacity int,data interface{}){
	defer ants.Release()
	p1, _ := ants.NewPool(capacity)
	syncCalculateSum := func() {
		f(data)
	}
	_ = p1.Submit(syncCalculateSum)
}

func Run(logger *zap.SugaredLogger,tailObj *tail.Tail,esObj *elasPkg.EsObj)  {
	capacity := 100
	// 处理并发送日志
	//SendLog()
	for {
		//logger.Warn("warning 日志...")
		//logger.Error("error 日志...")
		//logger.Infow("Infow自定义字段日志","name","zhang3","age",12)

		line, ok := <-tailObj.Lines //遍历chan，读取日志内容
		if !ok {
			fmt.Printf("tail file close reopen, filename:%s\n", inp.Inputlog)
			time.Sleep(500*time.Microsecond)
			continue
		}
		// tailf输出到终端
		//fmt.Println(line.Text)
		// tailf重写到logger
		logger.Info(line.Text)

		// tailf 推送日志导elasticsearch
		fmt.Println(esObj.IndexName)

		// 判断日志是否为json格式
		isJson := json.Valid([]byte(line.Text))
		//fmt.Println(isJson)
		f := func(data interface{}) {_ = esObj.PutData(data)}
		if isJson{
			threadPool(f,capacity,line.Text)
		}else {
			p := fmt.Sprintf(`{"logmsg":"%s"}`,strings.Replace(line.Text,"\r","",-1))
			//go f(p)
			threadPool(f,capacity,p)
		}
	}
}
