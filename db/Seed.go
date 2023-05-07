package db

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func SeedDB(db string) error {
	dbConnection, err := ConnectToDb(db)
	
	if err != nil {
		log.Panicln(err, "Error connecting to db", db)
	}
	productSalesRepo := NewProductSalesRepo(dbConnection)

	var products []Product
	for i := 0; i < 10; i++ {
		date := time.Date(2001, 12, rand.Intn(31)+1, rand.Intn(24), rand.Intn(60), rand.Intn(60), 0, time.Local)
		product := "Product " + string(i+65) // ASCII code for 'A' is 65
		region := "Region " + string(i+65)
		qty := uint32(rand.Int31n(100))
		cost := float32(rand.Intn(10000)) / 100.0
		tax := float32(rand.Intn(10000)) / 100.0
		products = append(products, Product{Date: date, Product: product, Region: region, Qty: qty, Cost: cost, Tax: tax})
	}

	for _, p := range products {
		_, err := productSalesRepo.CreateProduct(p)
		if err != nil {
			fmt.Println("Error occured", p, err)
			return err
		}
	}

	return nil
}
