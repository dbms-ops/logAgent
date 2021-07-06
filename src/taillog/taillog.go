package taillog

import (
	"fmt"
	"github.com/hpcloud/tail"
	"logAgent/src/kafka"
)

// 用于进行日志采集的模块

type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
}

func NewTailTask(path, topic string) (tailObj *TailTask) {
	tailObj = &TailTask{
		path:  path,
		topic: topic,
	}
	tailObj.init() // 根据路径打开对应的日志信息

	return

}

func (t *TailTask) init() {

	config := tail.Config{
		ReOpen:    true,                                 // 日志文件重新打开，日志文件会轮转
		Follow:    true,                                 // 日志文件follow
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 文件读取的位置
		MustExist: false,                                // 文件是否一定存在
		Poll:      false,                                //
	}
	var err error
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Println("tail file failed, err\n", err)
		return
	}
	go t.run()
	return
}

func (t *TailTask) ReadChan() <-chan *tail.Line {

	return t.instance.Lines
}

func (t *TailTask) run() {
	for {
		select {
		// 先将日志发送到一个通道中
		// kafka 那个包中有单独的goroutine 去日志数据发送到 kafka 中
		case line := <-t.instance.Lines: // 从 tailObj 通道中一行一行读取数据
			kafka.SentToChan(t.topic, line.Text) // 函数调用函数发送数据
		}
	}

}
