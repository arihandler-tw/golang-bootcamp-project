package broker

import (
	"encoding/json"
	"fmt"
	"gin-exercise/pkg/product/db"
	"github.com/Shopify/sarama"
	"log"
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

	partitionConsumer, err := consumer.ConsumePartition(ProductsTopic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	productsDatabase := db.NewProductsDatabase()
	requested := 0

	for {
		// blocks before select until there is a msg in at least one chan
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d: '%s': '%s'\n", msg.Offset, string(msg.Key), msg.Value)

			requested++

			var request ProductCreationRequest
			err := json.Unmarshal(msg.Value, &request)

			if err != nil {
				fmt.Printf("Failed to unmarshall message: %s", err)
				continue
			}

			storedProduct, err := productsDatabase.Store(request.Id, request.Price, request.Description)
			if err != nil {
				fmt.Printf("error during product store: %v", err)
			}
			fmt.Printf("product stored (%v)", storedProduct)
		}

		log.Printf("Store product requests: %d\n", requested)
	}
}
