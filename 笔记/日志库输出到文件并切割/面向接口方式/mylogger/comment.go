package mylogger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

type inputLevel uint8

const (
	UNKNOWN inputLevel = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)


const TIME_LAYOUT = "2006-01-02 15:04:05"
const TIME_LAYOUT_LOG = "2006-01-02_15:04:05"

const LOGPATH string = "."


func getInfo(skip int) (funcName, fileName string, lineNo int) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		fmt.Println("runtime.Caller() failed\n")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(file) // path.Base 获取路径最后的文件名
	funcName = strings.Split(funcName, ".")[1]
	return
}


// 输出到终端
func expConsole(mess string,lv string){
	funcName, fileName, lineNo := getInfo(3)
	now := time.Now().Format(TIME_LAYOUT)
	fmt.Printf("%v [%v] [%s:%s:%d] %s\n",now,lv,funcName,fileName,lineNo,mess)
}

// 输出到日志文件
func expLogFile(mess,lv,logpath string){
	now := time.Now().Format(TIME_LAYOUT)
	logfile := path.Join(logpath,"mylog.out")
	f,err := os.OpenFile(logfile,os.O_CREATE|os.O_APPEND|os.O_WRONLY,0644)
	defer f.Close()
	if err != nil{
		panic("日志文件创建失败")
		return
	}
	funcName, fileName, lineNo := getInfo(3)
	logmes := fmt.Sprintf("%v [%v] [%s:%s:%d] %s\n",now,lv,funcName,fileName,lineNo,mess)
	//fmt.Println(logmes)
	f.WriteString(logmes)
}

// 拷贝文件
func copyFile(src,dest string)(int64,error){
	srcfile, err := os.Open(src)
	if err != nil{
		fmt.Println(err)
		return 0, err
	}
	destfile , err := os.OpenFile(dest,os.O_RDWR|os.O_CREATE,os.ModePerm) // ModePerm默认0777
	if err !=nil {
		fmt.Println(err)
		return 0, err
	}
	defer srcfile.Close()
	defer destfile.Close()
	// 直接返回，在mian函数判断err省去一步if nil
	res ,err := io.Copy(destfile,srcfile)
	return res,err
}

// 按时间切割文件
func Cutfile(l *Logger) {
	//fmt.Println("0000000000")
	//fmt.Println(l)
	//fmt.Println(l.starttime)
	//fmt.Println(time.Now())
	if !l.starttime.Before(time.Now()){
		//fmt.Println("111111111")
		return
	}
	l.starttime = time.Now().Add(10*time.Second)
	//fmt.Println(l)
	//fmt.Println("2222222222222222222")
	if l.logpath == "" {
		l.logpath = LOGPATH
	}
	logfile := path.Join(l.logpath, "mylog.out") ////////////////////////////////
	//ticker := time.NewTicker(2*time.Second)
	//if Exists(logfile) {
	//	f, err := os.Create(logfile)
	//	if err != nil {
	//		panic("文件创建失败")
	//	}
	//	f.Close()
	////}
	//for{
	//	<-ticker.C // 输出时间，可以不打印直接写<-ticker.C
	now := time.Now().Format(TIME_LAYOUT_LOG)                        /////////////////////////
	dest := fmt.Sprintf("%v/%v_%v.log", l.logpath, "mylog.out", now) ////////////////////////////////
	dest2 := strings.Replace(dest, ":", "_", -1)
	_ , err := copyFile(logfile, dest2)
	if err != nil{
		fmt.Println("copy file faied")
		return
	}

	fff , err := os.OpenFile(logfile,os.O_WRONLY|os.O_TRUNC,os.ModePerm) // ModePerm默认0777
	if err !=nil {
		fmt.Println(err)
		return
	}
	fff.WriteString("")
	defer fff.Close()
	//}

}