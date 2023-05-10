package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ProductsDatabase() *gorm.DB {
	db, err := gorm.Open(postgres.Open(DSN()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&ProductEntity{})
	return db
}

func DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "5432", "postgres", "s3cr3t", "products")
}
