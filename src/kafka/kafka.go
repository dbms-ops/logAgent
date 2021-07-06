package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

// 专门往 kafka 写日志的模块

type logData struct {
	topic string
	data  string
}

var (
	client      sarama.SyncProducer // 声明一个全局的连接 kafka 生产者
	logDataChan chan *logData
)

// 1、初始化 kafka
func Init(address []string, maxSize int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // 发送完数据需要 leader 和 follow 都需要确认

	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选举出一个 partition
	config.Producer.Return.Successes = true                   // 成功交付消息，将在success channel 中返回

	// 连接 kafka
	client, err = sarama.NewSyncProducer(address, config) // 连接 kafka
	if err != nil {
		fmt.Println("producer closed,err:", err)
		return
	}
	fmt.Println("连接 kafka 成功")
	logDataChan = make(chan *logData, maxSize)
	// 开启后台的 goroutine 从通道中获取数据，发送到kafka
	go sendToKafka()
	return
}

// 给外部暴漏的一个函数,该函数只把日志数据发送到内部的channel
func SentToChan(topic, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}

// 向 kafka 中发送数据的函数
func sendToKafka() {

	for {
		select {
		case ld := <-logDataChan:
			msg := &sarama.ProducerMessage{} // 构造消息
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)
			// 发送 到 kafka
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send message failed. err ", err)
				return
			}
			fmt.Printf("pid:%v, offset:%v\n", pid, offset)
		default:
			time.Sleep(time.Millisecond * 50)
		}

	}
}
