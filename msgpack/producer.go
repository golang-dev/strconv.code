package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/vmihailenco/msgpack"
)

type Subject struct {
	ID     int
	Title  string
	Genres []string
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> <topic>\n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	topic := os.Args[2]

	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	producer, err := sarama.NewAsyncProducer([]string{broker}, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var enqueued, errors int
	doneCh := make(chan struct{})
	go func() {
		for {

			time.Sleep(1 * time.Second)

			subject := Subject{
				ID:     2,
				Title:  "千与千寻",
				Genres: []string{"剧情", "喜剧", "爱情", "战争"},
			}

			b, err := msgpack.Marshal(&subject)
			if err != nil {
				panic(err)
			}

			msg := &sarama.ProducerMessage{
				Topic: topic,
				Key:   sarama.StringEncoder(strconv.Itoa(subject.ID)),
				Value: sarama.StringEncoder(b),
			}
			select {
			case producer.Input() <- msg:
				enqueued++
				fmt.Printf("Produce message: %v\n", subject)
			case err := <-producer.Errors():
				errors++
				fmt.Println("Failed to produce message:", err)
			case <-signals:
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	log.Printf("Enqueued: %d; errors: %d\n", enqueued, errors)
}
