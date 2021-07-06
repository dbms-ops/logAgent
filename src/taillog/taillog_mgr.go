package taillog

import (
	"fmt"
	"logAgent/src/etcd"
	"time"
)

var tskMgr *tailLogMgr

// tailLog 对象的管理者,管理所有的 taillog 对象
type tailLogMgr struct {
	logEntry    []*etcd.LogEntry
	taskMap     map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

func Init(logEntryConf []*etcd.LogEntry) {
	tskMgr = &tailLogMgr{
		logEntry:    logEntryConf, // 把当前的日志收集项配置信息保存起来
		taskMap:     make(map[string]*TailTask, 32),
		newConfChan: make(chan []*etcd.LogEntry), // 无缓冲区的通道
	}

	for _, logEntry := range logEntryConf {
		NewTailTask(logEntry.Path, logEntry.Topic)
		// 包装为 发送日志的方法

	}
	go tskMgr.run()
}

// 监听自己的newConfChan,有了新的配置之后，进行对应的处理
// 1、配置新增
// 2、配置删除
// 3、配置变更
func (t *tailLogMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			fmt.Printf("this is new conf: %v \n", newConf)
		default:
			time.Sleep(time.Millisecond * 50)

		}
	}
}

// 通过函数向外部暴漏：tskMgr的newConfChan

func NewConfChan() chan<- []*etcd.LogEntry {
	return tskMgr.newConfChan
}
