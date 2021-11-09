package mynsq

import (
	"fmt"
	"github.com/nsqio/go-nsq"
)

var producer = InitProducer()

func ProduceSnatchMessage(message string)  {
	if err := producer.Publish("snatch", []byte(fmt.Sprintf(message))); err != nil {
		fmt.Println("Publish", err)
		panic(err)
	}
}

func ProduceOpenMessage(message string)  {
	if err := producer.Publish("open", []byte(fmt.Sprintf(message))); err != nil {
		fmt.Println("Publish", err)
		panic(err)
	}
}

func ProduceWalletMessage(message string)  {
	if err := producer.Publish("wallet", []byte(fmt.Sprintf(message))); err != nil {
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

//func main() {
//	producer, err := nsq.NewProducer("111.62.107.178:4150", nsq.NewConfig())
//	if err != nil {
//		fmt.Println("NewProducer", err)
//		panic(err)
//	}
//	//var message string
//	//fmt.Scanln(&message)
//	for {
//		if err := producer.Publish("test", []byte(fmt.Sprintf("hello"))); err != nil {
//			fmt.Println("Publish", err)
//			panic(err)
//		}
//	}
//}
