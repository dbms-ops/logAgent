package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Printf("connect etcd failed, error: %v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()
	// watch demo

	// 派一个哨兵，一直监视着某个 key:新增、修改、删除;返回一个channel
	ch := cli.Watch(context.Background(), "lixun")

	for wresp := range ch {
		for _, evt := range wresp.Events {
			fmt.Printf("type:%T key:%v value:%v", evt.Type, string(evt.Kv.Key), string(evt.Kv.Value))
		}
	}

}
