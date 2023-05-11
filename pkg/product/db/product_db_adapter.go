package db

import (
	"errors"
	"gin-exercise/pkg/product/model"
	"gin-exercise/pkg/util"
	"github.com/rs/xid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	db *gorm.DB
}

func NewProductsDatabase() *Repository {
	db, err := gorm.Open(postgres.Open(DSN()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&ProductEntity{})
	if err != nil {
		panic("failed to migrate database")
	}
	return &Repository{db: db}
}

func (d *Repository) Store(id *string, price float32, description string) (*model.Product, error) {
	if id == nil {
		id = util.GetPtr(generateNewId())
	}
	newProduct, newErr := model.NewProduct(*id, price, description, time.Now())
	if newErr != nil {
		return nil, newErr
	}
	productEntity := FromProduct(newProduct)
	tx := d.db.Create(&productEntity)
	return newProduct, tx.Error
}

func (d *Repository) Find(id string) (*model.Product, bool) {
	var prdEnt ProductEntity
	found := d.db.First(&prdEnt, "id = ?", id)
	product, _ := prdEnt.ToProduct()
	return product, !errors.Is(found.Error, gorm.ErrRecordNotFound)
}

func (d *Repository) Delete(id string) bool {
	tx := d.db.Delete(&ProductEntity{}, "id = ?", id)
	return tx.RowsAffected != 0
}

func (d *Repository) GetMany(amount int) ([]model.Product, error) {
	var entities []ProductEntity
	tx := d.db.Limit(amount).Find(&entities)
	if err := tx.Error; err != nil {
		return nil, err
	}

	var products []model.Product
	for _, entity := range entities {
		product, err := entity.ToProduct()
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
	}

	return products, nil
}

func generateNewId() string {
	return xid.New().String()
}
