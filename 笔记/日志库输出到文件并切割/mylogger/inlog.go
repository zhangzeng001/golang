package mylogger

import (
	"errors"
	"os"
	"strings"
	"time"
)

type logInterface interface {
	Debug()
	Trace()
	Info()
	Warning()
	Error()
	Fatal()
	cutfile()
}

type logger struct {
	logLevel inputLevel
	isConsole bool
	logpath string
	starttime time.Time
}

func NewlogObj(level string,con bool,logpath string)(logger,error){
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
		return logger{logLevel: UNKNOWN}, err
	}
	return logger{
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

func (l *logger) Debug(mess string){
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
func (l *logger) Trace(mess string){
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
func (l *logger) Info(mess string){
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
func (l *logger) Warning(mess string){
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
func (l *logger) Error(mess string){
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
func (l *logger) Fatal(mess string){
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
