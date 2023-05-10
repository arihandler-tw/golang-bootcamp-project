package product

import (
	"fmt"
)

func DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "5432", "postgres", "s3cr3t", "products")
}
