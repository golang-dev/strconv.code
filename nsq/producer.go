package main

import (
	"github.com/nsqio/go-nsq"
	"log"
	"math/rand"
	"time"
)

func main() {
	config := nsq.NewConfig()
	w, err := nsq.NewProducer("127.0.0.1:4150", config)

	if err != nil {
		log.Panic(err)
	}

	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	for {
		buf := make([]byte, 4)
		for i := 0; i < 4; i++ {
			buf[i] = chars[rand.Intn(len(chars))]
		}
		log.Printf("Pub: %s", buf)
		err = w.Publish("test", buf)
		if err != nil {
			log.Panic(err)
		}
		time.Sleep(time.Second * 1)
	}

	w.Stop()
}
