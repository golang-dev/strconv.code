package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/vmihailenco/msgpack"
)

type Subject struct {
	ID     int
	Title  string
	Genres []string
}

type Consumer struct {
	ready chan bool
}
var subject Subject

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		err := msgpack.Unmarshal(message.Value, &subject)
		if err != nil {
			panic(err)
		}
		log.Printf("Message claimed: key = %s, Subject(id=%d, title=%s, genres=%v)", string(message.Key),
			subject.ID, subject.Title, subject.Genres)
		session.MarkMessage(message, "")
	}

	return nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> <group> <topics..>\n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	group := os.Args[2]
	topics := os.Args[3:]

	version, err := sarama.ParseKafkaVersion("2.3.0")
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	config := sarama.NewConfig()
	config.Version = version
	consumer := Consumer{
		ready: make(chan bool, 0),
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup([]string{broker}, group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			if err := client.Consume(ctx, topics, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool, 0)
		}
	}()

	<-consumer.ready

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}
