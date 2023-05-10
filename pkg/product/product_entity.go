package product

import (
	"time"
)

type ProductEntity struct {
	ID          string `gorm:"primaryKey"`
	CreatedAt   time.Time
	Price       float32
	Description string
}

func FromProduct(product *Product) *ProductEntity {
	return &ProductEntity{
		ID:          product.ID,
		CreatedAt:   product.Creation,
		Price:       product.Price,
		Description: product.Description,
	}
}

func (p *ProductEntity) ToProduct() (*Product, error) {
	return NewProduct(p.ID, p.Price, p.Description, p.CreatedAt)
}
