package etcd

import (
	"github.com/coreos/etcd/clientv3"
	"time"
)

import (
	"context"
	"fmt"
)

func etcdPut() {

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
	cxt, cancel := context.WithTimeout(context.Background(), time.Second)
	value := `[
				{"path":"/data/redis.log","topic":"redis_log"},
				{"path":"/data/mongo.log","topic":"mongo_log"},
				{"path":"/data/mongo.log","topic":"mongo_log"}
			]`
	_, err = cli.Put(cxt, "/log_agent/collect_config", value)
	if err != nil {

		fmt.Println("put key failed, error:", err)
		return
	}
	cxt, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(cxt, "lixun")
	cancel()
	if err != nil {
		fmt.Printf("get from etcd faile,err \n", err)
		return
	}
	// 一次取出里面几乎所有的 key
	for _, ev := range resp.Kvs {
		fmt.Println("%s:%s\n", ev.Key, ev.Value)
	}

}
