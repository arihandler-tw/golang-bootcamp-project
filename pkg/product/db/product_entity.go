package db

import (
	"gin-exercise/pkg/product/model"
	"time"
)

type ProductEntity struct {
	ID          string `gorm:"primaryKey"`
	CreatedAt   time.Time
	Price       float32
	Description string
}

func FromProduct(product *model.Product) *ProductEntity {
	return &ProductEntity{
		ID:          product.ID,
		CreatedAt:   product.Creation,
		Price:       product.Price,
		Description: product.Description,
	}
}

func (p *ProductEntity) ToProduct() (*model.Product, error) {
	return model.NewProduct(p.ID, p.Price, p.Description, p.CreatedAt)
}
