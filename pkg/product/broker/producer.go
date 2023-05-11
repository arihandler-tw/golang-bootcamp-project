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

func (p ProductEventProducer) SendEvent(id *string, price float32, description string) {

	body, _ := json.Marshal(ProductCreationRequest{
		id:          id,
		price:       price,
		description: description,
	})

	msg := sarama.ProducerMessage{
		Topic: ProductsTopic,
		Key:   sarama.StringEncoder(xid.New().String()),
		Value: sarama.StringEncoder(body),
	}

	p.producer.Input() <- &msg
}
