package broker

import (
	"fmt"
	"gin-exercise/pkg/product/db"
	"gin-exercise/pkg/product/model"
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

	creationConsumer, err := consumer.ConsumePartition(ProductsCreationTopic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	deletionConsumer, err := consumer.ConsumePartition(ProductsDeletionTopic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := creationConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
		if err := deletionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	productsDatabase := db.NewProductsDatabase()

	for {
		// blocks before select until there is a msg in at least one chan
		select {
		case msg := <-creationConsumer.Messages():
			log.Printf("Received creation event (%v|%v): %v", string(msg.Key), msg.Offset, string(msg.Value))
			product, creationErr := handleCreationRequest(msg, productsDatabase)
			if creationErr != nil {
				log.Printf("[ERROR] creation of the product: %v", creationErr)
				continue
			}
			log.Printf("Product created %v\n", product)

		case msg := <-deletionConsumer.Messages():
			log.Printf("Received deletion event (%v|%v): %v", string(msg.Key), msg.Offset, string(msg.Value))
			deleted, deletionErr := handleDeletionRequest(msg, productsDatabase)
			if deletionErr != nil {
				log.Printf("[ERROR] deletion of the product: %v", deletionErr)
				continue
			}
			log.Printf("Product deleted %v\n", deleted)
		}
	}
}

func handleCreationRequest(msg *sarama.ConsumerMessage, database *db.Repository) (prd *model.Product, err error) {
	var request ProductCreationRequest
	err = Unmarshal(msg.Value, &request)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal: %w", err)
	}

	prd, err = database.Store(request.Id, request.Price, request.Description)
	if err != nil {
		return nil, fmt.Errorf("could not store product: %w", err)
	}
	return prd, nil
}

func handleDeletionRequest(msg *sarama.ConsumerMessage, database *db.Repository) (deleted bool, err error) {
	var request ProductDeletionRequest
	err = Unmarshal(msg.Value, &request)

	if err != nil {
		return false, fmt.Errorf("could not unmarshal: %w", err)
	}

	deleted = database.Delete(request.Id)
	return deleted, nil
}
