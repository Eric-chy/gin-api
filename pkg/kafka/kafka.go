//Package kafka 暂未进行封装，只简单写了生产者和消费者例子
package kafka

import (
	"fmt"
	"ginpro/config"
	"ginpro/pkg/helper/convert"
	"ginpro/pkg/helper/gjson"
	"ginpro/pkg/helper/gtime"
	"github.com/Shopify/sarama"
)

//Producer 生产者
func Producer(Msg string) {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	conf.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	conf.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	cfg := config.Conf.Kafka
	msg.Topic = cfg.Topic
	base := make(map[string]interface{})
	base["timestamp"] = gtime.GetMicroTime()
	base["level"] = "error"
	base["data"] = Msg
	json := gjson.JsonEncode(base)
	msg.Value = sarama.StringEncoder(json)
	// 连接kafka
	connect := convert.SplitAndTrim(cfg.Connect, ",")
	client, err := sarama.NewSyncProducer(connect, conf)
	if err != nil {
		return
	}
	defer client.Close()
	// 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}

//Consumer 消费者
func Consumer() {
	cfg := config.Conf.Kafka
	connect := convert.SplitAndTrim(cfg.Connect, ",")
	consumer, err := sarama.NewConsumer(connect, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	topic := cfg.Topic
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println(partitionList)
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(pc)
	}
}
