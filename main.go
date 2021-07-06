package main

import (
	"gopkg.in/ini.v1"
	"logAgent/src/conf"
	"logAgent/src/etcd"
	"logAgent/src/kafka"
	"logAgent/src/taillog"
	"logAgent/src/utils"
	"sync"
	"time"
)

import (
	"fmt"
)

var (
	cafg = new(conf.AppConf)
)

// 日志采集模块的入口

func main() {
	// 0、加载配置文件
	err := ini.MapTo(cafg, "./conf/config.go")
	if err != nil {
		fmt.Printf("load ini failed, err: %v \n", err)
		return
	}

	// 1、初始化 kafka 连接
	err = kafka.Init([]string{cafg.KafkaConf.Address}, cafg.KafkaConf.ChanMaxSize)
	if err != nil {
		fmt.Printf("init kafka error:%v\n", err)
		return
	}
	// 2、初始化 etcd 连接

	err = etcd.Init(cafg.EtcdConf.Address, time.Duration(cafg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("init etcd failed, err: %v", err)
		return
	}
	fmt.Println("init etcd successful")
	// 为了实现 logagent 都拉取自己独有的配置，需要通过自己的IP地址进行区分
	ipStr, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	etcdConfKey := fmt.Sprintf(cafg.EtcdConf.Key, ipStr)
	// 2.1 从 etcd 中获取日志的配置信息
	logEntryConf, err := etcd.GetConf(etcdConfKey)
	if err != nil {
		fmt.Printf("get conf from etcd failed, err: %v", err)
		return
	}
	fmt.Printf("get conf from etcd successed, %v \n", logEntryConf)
	// 2.2 使用哨兵监控日志收集项目的变化，有变化及时通知 LogAgent,实现热加载

	for index, value := range logEntryConf {
		fmt.Printf("index:%v value:%v \n", index, value)
	}

	// 3 收集日志、发送到 kafka 中
	// 3.1 循环每一个收集项，创建 TopicObj
	taillog.Init(logEntryConf)
	// Init 在 NewConfChan 之前，防止空指针
	newConfChan := taillog.NewConfChan() // 从taillog 包中获取对外暴漏的通道
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(cafg.EtcdConf.Key, newConfChan) // 哨兵发现新的配置可以通知上面的通道
	wg.Wait()
	// 3.2 发送到 KafkaZ

}
