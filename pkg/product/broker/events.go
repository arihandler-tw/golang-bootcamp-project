package broker

import (
	"encoding/json"
)

const (
	ProductsCreationTopic = "create-products"
	ProductsDeletionTopic = "delete-products"
)

type ProductCreationRequest struct {
	Id          *string
	Price       float32
	Description string
}

type ProductDeletionRequest struct {
	Id string
}

type ProductOperationRequest interface {
	ProductCreationRequest | ProductDeletionRequest
}

func Unmarshal[P ProductOperationRequest](msg []byte, req *P) (err error) {
	return json.Unmarshal(msg, &req)
}
