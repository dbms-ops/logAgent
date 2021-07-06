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

		// 包装为 发送日志的方法
		tailObj := NewTailTask(logEntry.Path, logEntry.Topic)
		mk := fmt.Sprintf("%s_%s", logEntry.Path, logEntry.Topic)
		tskMgr.taskMap[mk] = tailObj

	}
	go tskMgr.run()
}

// 监听自己的newConfChan,有了新的配置之后，进行对应的处理

func (t *tailLogMgr) run() {
	for {
		select {

		case newConf := <-t.newConfChan:
			// 1、配置新增
			for _, conf := range newConf {
				mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				_, ok := t.taskMap[mk]
				if ok {
					// 如果存在，则跳过
					continue
				} else {
					// 如果不存在，则继续
					tailObj := NewTailTask(conf.Path, conf.Topic)
					t.taskMap[mk] = tailObj
				}
			}
			// 找出原来的 t.LogEntry 有，但是 newConf中没有的，需要删除掉
			for _, c1 := range t.logEntry {
				isDelete := true
				for _, c2 := range newConf {
					if c2.Path == c1.Path && c2.Topic == c1.Topic {
						continue
					}
				}
				if isDelete {
					// 把 c1 对应的tailObj 停止掉
					mk := fmt.Sprintf("%s_%s", c1.Path, c1.Topic)
					// 调用 ctx.cancel 结束 t.run()函数
					t.taskMap[mk].cancelFunc()
				}
			}

			// 2、配置删除
			// 3、配置变更

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
