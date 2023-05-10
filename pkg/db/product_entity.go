package db

import (
	"gin-exercise/pkg/product"
	"time"
)

type ProductEntity struct {
	ID          string `gorm:"primaryKey"`
	CreatedAt   time.Time
	Price       float32
	Description string
}

func ProductEntityFromProduct(product *product.Product) *ProductEntity {
	return &ProductEntity{
		ID:          product.ID,
		CreatedAt:   product.Creation,
		Price:       product.Price,
		Description: product.Description,
	}
}

func (p *ProductEntity) ToProduct() (*product.Product, error) {
	return product.NewProduct(p.ID, p.Price, p.Description, p.CreatedAt)
}
