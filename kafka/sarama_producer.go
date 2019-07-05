package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
)

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

	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	var enqueued, errors int
	doneCh := make(chan struct{})
	go func() {
		for {

			time.Sleep(1 * time.Second)

			buf := make([]byte, 4)
			for i := 0; i < 4; i++ {
				buf[i] = chars[rand.Intn(len(chars))]
			}

			strTime := strconv.Itoa(int(time.Now().Unix()))
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Key:   sarama.StringEncoder(strTime),
				Value: sarama.StringEncoder(buf),
			}
			select {
			case producer.Input() <- msg:
				enqueued++
				fmt.Printf("Produce message: %s\n", buf)
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
