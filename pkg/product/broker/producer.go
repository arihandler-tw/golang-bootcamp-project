package broker

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/rs/xid"
)

type ProductEventProducer struct {
	producer sarama.AsyncProducer
}

func NewEventProducer() (*ProductEventProducer, error) {
	config := sarama.NewConfig()
	producer, err := sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		return nil, err
	}
	return &ProductEventProducer{
		producer: producer,
	}, nil
}

func (p ProductEventProducer) SendCreationRequest(id *string, price float32, description string) {

	body, err := json.Marshal(ProductCreationRequest{
		Id:          id,
		Price:       price,
		Description: description,
	})

	if err != nil {
		panic(err)
	}

	p.inputMessage(body, ProductsCreationTopic)
}

func (p ProductEventProducer) SendDeletionRequest(id string) {
	body, err := json.Marshal(ProductDeletionRequest{
		Id: id,
	})

	if err != nil {
		panic(err)
	}

	p.inputMessage(body, ProductsDeletionTopic)
}

func (p ProductEventProducer) inputMessage(body []byte, topic string) {
	msg := sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(xid.New().String()),
		Value: sarama.StringEncoder(body),
	}

	p.producer.Input() <- &msg
}
