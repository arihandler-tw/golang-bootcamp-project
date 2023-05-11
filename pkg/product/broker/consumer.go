package broker

import (
	"encoding/json"
	"fmt"
	"gin-exercise/pkg/product/db"
	"gin-exercise/pkg/server"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
)

func Consumer() {
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, sarama.NewConfig())
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	topic := "products"
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	fmt.Printf("Creating database:...")
	r := db.NewProductsDatabase()
	fmt.Printf("Database created:...")
	consumed := 0
	notKilled := true

	for notKilled {
		// blocks before select until there is a msg in at least one chan
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d: '%s': '%s'\n", msg.Offset, string(msg.Key), msg.Value)

			consumed++

			k := string(msg.Key)
			var p server.ProdReq
			err := json.Unmarshal(msg.Value, &p)

			if err != nil {
				fmt.Printf("Failed to unmarshall message: %s", err)
				continue
			}

			r.Store(&k, p.Price, p.Description)

		case <-signals:
			notKilled = false
		}
	}

	log.Printf("Consumed: %d\n", consumed)
}
