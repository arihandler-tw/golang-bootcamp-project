package product

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type DBRepository struct {
	db *gorm.DB
}

func NewProductsDatabase() *DBRepository {
	db, err := gorm.Open(postgres.Open(DSN()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&ProductEntity{})
	return &DBRepository{db: db}
}

func (d *DBRepository) Store(id *string, price float32, description string) (*Product, error) {
	newProduct, newErr := NewProduct(*id, price, description, time.Now())
	if newErr != nil {
		return nil, newErr
	}
	productEntity := FromProduct(newProduct)
	tx := d.db.Create(&productEntity)
	return newProduct, tx.Error
}

func (d *DBRepository) Find(id string) (*Product, bool) {
	var prdEnt ProductEntity
	found := d.db.First(&prdEnt, "id = ?", id)
	product, _ := prdEnt.ToProduct()
	return product, !errors.Is(found.Error, gorm.ErrRecordNotFound)
}

func (d *DBRepository) Delete(id string) bool {
	tx := d.db.Delete(&ProductEntity{}, "id = ?", id)
	return tx.RowsAffected != 0
}

func (d *DBRepository) GetMany(amount int) ([]Product, error) {
	var entities []ProductEntity
	tx := d.db.Limit(amount).Find(&entities)
	if err := tx.Error; err != nil {
		return nil, err
	}

	var products []Product
	for _, entity := range entities {
		product, err := entity.ToProduct()
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
	}

	return products, nil
}
