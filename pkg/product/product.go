package product

import (
	"errors"
	"time"
	"unicode"
)

type Product struct {
	ID          string
	Price       float32
	Creation    time.Time
	Description string
}

func NewProduct(id string, price float32, description string, creation time.Time) (p *Product, e error) {
	if !isASCII(id) {
		return nil, errors.New("ID should contain only ASCII characters")
	}
	if len(description) > 50 {
		return nil, errors.New("description should be less than 50 characters long")
	}

	return &Product{
		ID:          id,
		Price:       price,
		Creation:    creation,
		Description: description,
	}, nil
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
