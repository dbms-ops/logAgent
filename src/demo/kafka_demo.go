package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // 发送完数据需要 leader 和 follow 都需要确认

	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选举出一个 partition
	config.Producer.Return.Successes = true                   // 成功交付消息，将在success channel 中返回

	msg := &sarama.ProducerMessage{} // 构造消息
	msg.Topic = "web_log"
	msg.Value = sarama.StringEncoder("this is a test log")

	client, err := sarama.NewSyncProducer([]string{"192.168.3.134:9092"}, config) // 连接 kafka
	if err != nil {
		fmt.Println("producer closed,err:", err)
		return
	}
	fmt.Println("连接 kafka 成功")
	defer client.Close()
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send message failed. err ", err)
		return
	}
	fmt.Printf("pid:%v, offset:%v\n", pid, offset)
}
