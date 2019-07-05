package main

import (
	"log"
	"sync"

	"github.com/nsqio/go-nsq"
)

func main() {

	wg := &sync.WaitGroup{}
	wg.Add(1000)

	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer("test", "ch", config)
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("Got a message: %s", message.Body)
		wg.Done()
		return nil
	}))
	err := q.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		log.Panic(err)
	}
	wg.Wait()

}
