# go操作kafka

## kafka工作原理

https://www.cnblogs.com/sujing/p/10960832.html

https://developer.ibm.com/zh/articles/os-cn-kafka/

kafka集群的架构

1.  **Producer**：Producer即生产者，消息的产生者，是消息的入口。
2.  **kafka cluster**：
      　　**Broker**：Broker是kafka实例，每个服务器上有一个或多个kafka的实例，我们姑且认为每个broker对应一台服务器。每个kafka集群内的broker都有一个**不重复**的编号，如图中的broker-0、broker-1等……
      　　**Topic**：消息的主题，可以理解为消息的分类，kafka的数据就保存在topic。在每个broker上都可以创建多个topic。
      　　**Partition**：Topic的分区，每个topic可以有多个分区，分区的作用是做负载，提高kafka的吞吐量。同一个topic在不同的分区的数据是不重复的，partition的表现形式就是一个一个的文件夹！
      　　**Replication**:每一个分区都有多个副本，副本的作用是做备胎。当主分区（Leader）故障的时候会选择一个备胎（Follower）上位，成为Leader。在kafka中默认副本的最大数量是10个，且副本的数量不能大于Broker的数量，follower和leader绝对是在不同的机器，同一机器对同一个分区也只可能存放一个副本（包括自己）。
   3. **Message**：每一条发送的消息主体。
      　**Consumer**：消费者，即消息的消费方，是消息的出口。
         　**Consumer Group**：我们可以将多个消费组组成一个消费者组，在kafka的设计中同一个分区的数据只能被消费者组中的某一个消费者消费。同一个消费者组的消费者可以消费同一个topic的不同分区的数据，这也是为了提高kafka的吞吐量！
         　**Zookeeper**：kafka集群依赖zookeeper来保存集群的的元信息，来保证系统的可用性。 

生产者向kafka发送数据的流程

![1613486506624](D:\go笔记\notes\笔记\kafka.assets\1613486506624.png)

kafka选择分区的3种模式

1. 指定往哪个分区写
2. 指定key，kafka根据key做hash然后决定写哪个分区
3. 轮询方式

生产者往kafka发送数据的3种模式

1.   `0` 把数据发送leader就返回成功，效率最高、安全性最低
2.   `1`把数据发送给leader，等待leader回ACK
3.   `all` 把数据发送给leader，确保follwer从leader拉去数据恢复ack给leader，leader再恢复ack，安全性最高

## 安装

windows使用1.19.0

`require github.com/Shopify/sarama v1.19.0`

```bash
go get github.com/Shopify/sarama
```

## 注意事项

`sarama` v1.20之后的版本加入了`zstd`压缩算法，需要用到cgo，在Windows平台编译时会提示类似如下错误：

```bash
# github.com/DataDog/zstd
exec: "gcc":executable file not found in %PATH%
```

所以在Windows平台请使用v1.19版本的sarama。



## go作为生产者

 ```go
package main

import (
	"fmt"

	"github.com/Shopify/sarama"
)

// 基于sarama第三方库开发的kafka client

func main() {
	config := sarama.NewConfig()
	// 生产者往kafka发送数据的种模式  WaitForAll=-1  NoResponse=0 WaitForLocal=1
	config.Producer.RequiredAcks = sarama.WaitForAll
	// kafka选择分区的3种模式
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 成功交付的消息将在success channel返回
	config.Producer.Return.Successes = true

	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	// 发送的topic名称
	msg.Topic = "web_log"
	msg.Value = sarama.StringEncoder("this is a test log")
	// 连接kafka []string{"10.0.0.11:9092","10.0.0.10:9092","10.0.0.12:9092"}
	client, err := sarama.NewSyncProducer([]string{"10.0.0.11:9092"}, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	defer client.Close()
	// 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}
	fmt.Println("发送成功...")
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}
 ```



## 消费者

```go
package main

import (
	"fmt"

	"github.com/Shopify/sarama"
)

// kafka consumer  消费消息

func main() {
	// 定义consumer OBJ
	consumer, err := sarama.NewConsumer([]string{"10.0.0.11:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	// 根据topic取到所有的分区
	partitionList, err := consumer.Partitions("web_log")
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println("partitionlist--",partitionList)
	// 遍历所有的分区
	for partition := range partitionList {
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		//go func(sarama.PartitionConsumer) {
		//	for msg := range pc.Messages() {
		//		fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
		//	}
		//}(pc)
		// 同步消费消息
		for msg := range pc.Messages() {
			fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
		}
	}
}
```



## kafka操作

**启动正确的zookeeper集群**

```
bin/zookeeper-server-start.sh config/zookeeper.properties
```

不同的broker-0编号，数据目录

https://www.orchome.com/6

**启动kafka**

```
bin/kafka-server-start.sh config/server.properties
```

**创建topic**

 kafka版本 >= 2.2 

```go
bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic test
```

 kafka版本 < 2.2 

```go
bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test
```

**列出topic**

```go
bin/kafka-topics.sh --list --bootstrap-server localhost:9092
```

**消费消息**

```go
bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic web_log --from-beginning
```

**发送消息**

```go
bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic web_log --from-beginning
```

