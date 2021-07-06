package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

// tailf 用法示例
func main() {
	fileName := "./my.log"
	config := tail.Config{
		ReOpen:    true,                                 // 日志文件重新打开，日志文件会轮转
		Follow:    true,                                 // 日志文件follow
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 文件读取的位置
		MustExist: false,                                // 文件是否一定存在
		Poll:      false,                                //
	}
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tail file failed, err\n", err)
		return
	}
	var (
		line *tail.Line
		ok   bool
	)
	for true {
		line, ok = <-tails.Lines
		if !ok {
			fmt.Printf("tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(time.Second)
		}
		fmt.Println("line: ", line.Text)
	}

}
