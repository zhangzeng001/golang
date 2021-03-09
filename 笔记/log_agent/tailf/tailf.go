package tailf

import (
	"fmt"
	"github.com/hpcloud/tail"
)

var (
	line *tail.Line
	ok   bool
)

func InitTail(logfile string) (*tail.Tail,error) {
	fileName := logfile
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return nil,err
	}

	return tails,nil
}
