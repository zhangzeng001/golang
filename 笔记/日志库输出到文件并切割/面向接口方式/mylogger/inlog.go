package mylogger

import (
	"errors"
	"os"
	"strings"
	"time"
)

type LogInterface interface {
	Debug(mess string)
	Trace(mess string)
	Info(mess string)
	Warning(mess string)
	Error(mess string)
	Fatal(mess string)
}

type Logger struct {
	logLevel inputLevel
	isConsole bool
	logpath string
	starttime time.Time
}


func NewlogObj(level string,con bool,logpath string)(*Logger,error){
	lv := strings.ToUpper(level)
	var res inputLevel
	switch lv {
	case "DEBUG":
		res = DEBUG
	case "TRACE":
		res = TRACE
	case "INFO":
		res = INFO
	case "WARNING":
		res = WARNING
	case "ERROR":
		res = ERROR
	case "FATAL":
		res = FATAL
	default:
		err := errors.New("无效的日志级别")
		return &Logger{logLevel: UNKNOWN}, err
	}
	return &Logger{
		logLevel: res,
		isConsole: con,
		logpath: logpath,
		starttime: time.Now(),
	},nil
}

// 判断所给路径文件/文件夹是否存在1
func Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (l *Logger) Debug(mess string){
	//Cutfile(&l)
	if l.logLevel < DEBUG{
		if l.isConsole{
			expConsole(mess,"DEBUG")
		}else {
			if l.logpath == ""{
				l.logpath = LOGPATH
			}
			Cutfile(l)
			expLogFile(mess,"DEBUG",l.logpath)
		}
	}
}
func (l *Logger) Trace(mess string){
	//Cutfile(&l)
	if l.logLevel < TRACE{
		if l.isConsole{
			expConsole(mess,"TRACE")
		}else {
			if l.logpath == ""{
				l.logpath = LOGPATH
			}
			Cutfile(l)
			expLogFile(mess,"TRACE",l.logpath)
		}
	}
}
func (l *Logger) Info(mess string){
	//Cutfile(&l)
	if l.logLevel < INFO{
		if l.isConsole{
			expConsole(mess,"INFO")
		}else {
			if l.logpath == ""{
				l.logpath = LOGPATH
			}
			Cutfile(l)
			expLogFile(mess,"INFO",l.logpath)
		}
	}
}
func (l *Logger) Warning(mess string){
	//Cutfile(&l)
	if l.logLevel < WARNING{
		if l.isConsole{
			expConsole(mess,"WARNING")
		}else {
			if l.logpath == ""{
				l.logpath = LOGPATH
			}
			Cutfile(l)
			expLogFile(mess,"WARNING",l.logpath)
		}
	}
}
func (l *Logger) Error(mess string){
	if l.logLevel < ERROR{
		if l.isConsole{
			expConsole(mess,"ERROR")
		}else {
			if l.logpath == ""{
				l.logpath = LOGPATH
			}
			Cutfile(l)
			expLogFile(mess,"ERROR",l.logpath)
		}
	}
}
func (l *Logger) Fatal(mess string){
	//Cutfile(&l)
	if l.logLevel < FATAL{
		if l.isConsole{
			expConsole(mess,"FATAL")
		}else {
			if l.logpath == ""{
				l.logpath = LOGPATH
			}
			Cutfile(l)
			expLogFile(mess,"FATAL",l.logpath)
		}
	}
}
