package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"

	"time"
)

var (
	cli *clientv3.Client
)

// 需要收集的日志的配置信息
type LogEntry struct {
	Path  string `json:"path"`  // 1、日志存放的路径
	Topic string `json:"topic"` // 2、日志需要发送的 topic
}

func Init(addr string, timeout time.Duration) (err error) {

	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: timeout,
	})
	if err != nil {
		fmt.Printf("connect etcd failed, error: %v\n", err)
		return
	}
	return
}

// 从 etcd 中根据 Key 获取配置项目
func GetConf(key string) (logEntryConf []*LogEntry, err error) {
	cxt, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(cxt, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd faile,err %s \n", err)
		return
	}
	// 将返回的 json 解压到 value 里面
	for _, env := range resp.Kvs {
		err = json.Unmarshal(env.Value, &logEntryConf)
		if err != nil {
			fmt.Println("unmarshall key value, err ", err)
			return
		}
	}

	return
}

func WatchConf(key string, newConfCh chan<- []*LogEntry) {
	ch := cli.Watch(context.Background(), key)
	for resp := range ch {
		for _, evt := range resp.Events {
			// 通知对应的接收者: taillog.tskMgr
			// 1、判断操作的类型
			var newConf []*LogEntry

			if evt.Type != clientv3.EventTypeDelete {
				// 如果是删除操作，手动传递一个空的 slice，否则解析会出现问题
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					fmt.Printf("unmarshall failed, err:%v\n", err)
					continue
				}

			}

			fmt.Printf("get new conf:%v \n", newConf)
			newConfCh <- newConf
			// 2、
		}
	}
}
