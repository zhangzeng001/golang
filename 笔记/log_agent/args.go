package main

import (
	"flag"
	"os"
)

// Argsconfig 处理命令行输出配置文件路径
func Argsconfig(){
	flag.Parse()
	if printHelp {
		flag.Usage()
		os.Exit(0)
	}
	if configName == "unknown"{
		//fmt.Println("USAGE:",os.Args[0],"--help")
		flag.Usage()
		os.Exit(1)
	}
}


