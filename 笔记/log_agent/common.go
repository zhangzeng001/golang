package main

import (
	"awesomeProject/log_agent/conf"
	"awesomeProject/log_agent/elasPkg"
	"awesomeProject/log_agent/tailf"
	"flag"
	"fmt"
	"github.com/hpcloud/tail"
	"go.uber.org/zap"
)

// 命令行指定配置文件
var (
	configName  string
	printHelp    bool
	logger *zap.SugaredLogger
	tailObj *tail.Tail
	esObj *elasPkg.EsObj
)

// ini 配置文件项目标题
const (
	inputlog = "inputlog"
	kafka = "kafka"
	logs = "logs"
	esSectionName = "elasticsearch"
)

// kafka 配置文件模块
type kafkaConf struct {
	Host string  `ini:"host"`
	Topic string `ini:"topic"`
}

// inputlog 配置文件模块
type inputLogConf struct {
	Inputlog string `ini:"inputlog"`
}

// logs 配置文件模块
type logsConf struct {
	Logdir string    `ini:"logdir"`
	Log_file string  `ini:"log_file"`
	Log_error string `ini:"log_error"`
	IsConsole bool   `ini:"isConsole"`
	Loglevel string  `ini:"loglevel"`
	MaxSize int      `ini:"MaxSize"`
	MaxBackups int   `ini:"MaxBackups"`
	MaxAge int       `ini:"MaxAge"`
	Compress bool    `ini:"Compress"`
}

// es 配置文件模块
type esConf struct {
	Host string    `ini:"host"`
	Index string   `ini:"index"`
}

// 配置文件解析变量
var (
	inp = new(inputLogConf)
	kaf = new(kafkaConf)
	es = new(esConf)
	logsSec = new(logsConf)
)

// analysis 解析ini配置文件
func Analysis(iniPath string,section string,confMap interface{})error{
	err := conf.Loadconf(iniPath,section,confMap)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}

func init()  {
	flag.BoolVar(&printHelp, "h", false, "Print this help")
	flag.StringVar(&configName,"config","unknown","--config /dir/config.ini")

	// 处理命令行参数
	Argsconfig()

	// inputlog
	_ = Analysis(configName,inputlog,inp)

	// kafka 配置解析
	_ = Analysis(configName,kafka,kaf)

	// logs 配置解析
	_ = Analysis(configName,logs,logsSec)

	// es 配置解析
	_ = Analysis(configName,esSectionName,es)

	// 初始化日志模块
	Loglumberjack := map[string]interface{}{
		"MaxSize": logsSec.MaxSize,         // 每个日志文件保存的大小 单位:M
		"MaxBackups": logsSec.MaxBackups,      // 日志文件最多保存多少个备份
		"MaxAge": logsSec.MaxAge,          // 日志文件最多保存多少天
		"Compress": logsSec.Compress,      // 是否压缩
	}
	logobj,_ := NewlogObj(logsSec.Loglevel,logsSec.IsConsole,logsSec.Log_file,Loglumberjack)
	logger = logobj.InitLogger()

	// 初始化tail
	var err error
	tailObj,_ = tailf.InitTail(inp.Inputlog)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 初始化es
	esObj = elasPkg.NewesObj(es.Host,es.Index)
}