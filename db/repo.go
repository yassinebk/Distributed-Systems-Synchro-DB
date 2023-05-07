package db

import (
	"log"

	"gorm.io/gorm"
)

type IProductSalesRepo interface {
	FindAll() []Product
	FindOne(int) *Product
	CreateProduct(Product) (Product, error)
	UpdateProduct(Product) (Product, error)
	DeleteProduct(int) (bool, error)
}

type ProductSalesRepo struct {
	db *gorm.DB
}

func NewProductSalesRepo(db *gorm.DB) ProductSalesRepo {
	return ProductSalesRepo{db: db}
}

func (p *ProductSalesRepo) FindAll() []Product {
	var products []Product
	res := p.db.Find(&products)

	if res.Error != nil {
		log.Panicln("[-] p.findAll - Error retrieving products sales from database")
	}

	return products
}

func (p *ProductSalesRepo) FindOne(id int) *Product {
	var product Product

	res := p.db.First(&product, id)

	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil
		}
		log.Panicln("[-] p.findOne - Error retrieving record for ", id, res.Error)
	}

	return &product
}

func (p *ProductSalesRepo) CreateProduct(newProduct Product) (Product, error) {
	res := p.db.Create(&newProduct)

	if res.Error != nil {
		return Product{}, res.Error
	}

	return newProduct, nil
}

func (p *ProductSalesRepo) UpdateProduct(updatedProduct Product) (Product, error) {
	res := p.db.Save(updatedProduct)

	if res.Error != nil {
		return Product{}, res.Error
	}

	return updatedProduct, nil
}

func (p *ProductSalesRepo) DeleteProduct(productId int) (bool, error) {

	res := p.db.Delete(&Product{}, productId)

	if res.Error != nil {
		return false, res.Error
	}

	return true, nil
}
