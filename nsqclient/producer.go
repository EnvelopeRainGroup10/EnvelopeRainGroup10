package nsqclient

import (
	"fmt"
	"github.com/nsqio/go-nsq"
)

var Producer = InitProducer()

func ProduceMessage(topic string, message string) {
	if err := Producer.Publish(topic, []byte(fmt.Sprintf(message))); err != nil {
		fmt.Println("Publish", err)
		panic(err)
	}
}

func InitProducer() *nsq.Producer {
	producer, err := nsq.NewProducer("111.62.107.178:4150", nsq.NewConfig())
	if err != nil {
		fmt.Println("NewProducer", err)
		panic(err)
	}
	return producer
}
