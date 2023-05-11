package product

import "gin-exercise/pkg/product/model"

var (
	IdTakenError = "id already registered"
)

type RepositoryInterface interface {
	Store(*string, float32, string) (*model.Product, error)
	Find(string) (*model.Product, bool)
	Delete(string) bool
	GetMany(int) ([]model.Product, error)
}
