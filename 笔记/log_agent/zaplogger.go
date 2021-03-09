package main
// 参考https://studygolang.com/articles/25044?fr=sidebar
import (
	"errors"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

// zaplogger 日志库相关结构体
type zaplogger struct {
	isConsole bool            // 命令行输出
	level    zapcore.Level    // 日志级别
	logpath  string           // 日志路劲
	Loglumberjack map[string]interface{}
	//MaxSize  uint32          // 每个日志文件保存的大小 单位:M
	//MaxBackups uint32        // 文件最多保存多少个
	//MaxAge   uint32          // 日志文件最多保存多少天
	//Compress  bool           // 是否压缩
	//sugarLogger *zap.SugaredLogger   // 定义全局*zap.SugaredLogger类型 sugarLogger 变量
}

// NewlogObj 构建函数
func NewlogObj(level string,isConsole bool,logpath string,Loglumberjack map[string]interface{})(*zaplogger,error){
	lv := strings.ToUpper(level)
	var res zapcore.Level
	switch lv {
	case "DEBUG":
		res = zapcore.DebugLevel
	case "INFO":
		res = zapcore.InfoLevel
	case "WARNING":
		res = zapcore.WarnLevel
	case "ERROR":
		res = zapcore.ErrorLevel
	case "PANIC":
		res = zapcore.PanicLevel
	default:
		err := errors.New("无效的日志级别")
		return nil, err
	}
	//类型判断和默认值设置
	MaxSize := Loglumberjack["MaxSize"]
	newMaxSize,_ := IntReflect(1024,MaxSize)
	MaxBackups := Loglumberjack["MaxBackups"]
	newMaxBackups,_ := IntReflect(0,MaxBackups)
	MaxAge := Loglumberjack["MaxAge"]
	newMaxAge,_ := IntReflect(15,MaxAge)
	Compress := Loglumberjack["Compress"]
	newCompress,_ := BoolReflect(false,Compress)

	Loglumberjack["MaxSize"] = newMaxSize
	Loglumberjack["MaxBackups"] = newMaxBackups
	Loglumberjack["MaxAge"] = newMaxAge
	Loglumberjack["Compress"] = newCompress

	return &zaplogger{
		level: res,
		isConsole: isConsole,
		logpath: logpath,
		Loglumberjack:Loglumberjack,
	},nil
}

// IntReflect 处理反射int类型信息
func IntReflect(def int,input interface{})(res int,err error){
	v := reflect.ValueOf(input)
	// 判断传入参数是否为空
	if !v.IsValid(){
		res = def
		return def,nil
	} else { // 非空
		k := v.Kind()

		switch k {
		case reflect.Int:
			return int(v.Int()),nil
		default:
			log.Fatal(input,", 输入类型错误")
			return 0,nil
		}
	}
}

// BoolReflect 处理bool反射值
func BoolReflect(def bool,input interface{})(res bool,err error){
	v := reflect.ValueOf(input)
	// 判断传入参数是否为空
	if !v.IsValid(){
		res = def
		return res,nil
	} else { // 非空
		k := v.Kind()

		switch k {
		case reflect.Bool:
			return v.Bool(),nil
		default:
			log.Fatal(input,", 输入类型错误")
			return false,nil
		}
	}
}

//////////////////////////////////     zap相关       ///////////////////////////////////////
// getEncoder NewConsoleEncoder方式输出
func getEncoder(isConsole bool) zapcore.Encoder {
	// 非json输出
	encoderConfig := zap.NewProductionEncoderConfig()
	// 时间编码
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	// 日志级别大写
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if isConsole{
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriter 构建 writeSyncer
func getLogWriter(logpath string,isConsole bool,Loglumberjack map[string]interface{}) zapcore.WriteSyncer {
	// 切割日志
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logpath,    // 日志文件路径
		MaxSize:    Loglumberjack["MaxSize"].(int),        // 每个日志文件保存的大小 单位:M
		MaxBackups: Loglumberjack["MaxBackups"].(int),     // 文件最多保存多少天
		MaxAge:     Loglumberjack["MaxAge"].(int),         // 日志文件最多保存多少个备份
		Compress:   Loglumberjack["Compress"].(bool),       // 是否压缩
	}
	// 得到一个文件句柄
	//file, _ := os.Create(logpath)

	// 终端打印
	if isConsole{
		return zapcore.AddSync(os.Stdout)
	}
	return zapcore.AddSync(lumberJackLogger)
}

// InitLogger 初始化日志对象
func (z *zaplogger)InitLogger() *zap.SugaredLogger {
	// 构建New(core)
	writeSyncer := getLogWriter(z.logpath,z.isConsole,z.Loglumberjack)
	encoder := getEncoder(z.isConsole)
	core := zapcore.NewCore(encoder, writeSyncer, z.level)
	//core := zapcore.NewCore(encoder,  zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), z.level)
	// zap.AddCaller() 打印函数
	logger := zap.New(core,zap.AddCaller())
	// 获取sugar对象
	sugarLogger := logger.Sugar()
	return sugarLogger
}
