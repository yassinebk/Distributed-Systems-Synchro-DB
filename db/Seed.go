package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func SeedDB(db string) error {
	dbConnection, err := ConnectToDb(db)

	if err != nil {
		log.Panicln(err, "Error connecting to db", db)
	}
	productSalesRepo := NewProductSalesRepo(dbConnection)

	var products []Product

	data, err := os.ReadFile("init-db.json")

	if err != nil {
		log.Panicln("[-] Error while seeding database. File open error", err)
	}

	err = json.Unmarshal(data, &products)

	if err != nil {
		log.Panicln("[-] Error unmarshelling data")
	}

	for _, p := range products {
		_, err := productSalesRepo.CreateProduct(p)
		if err != nil {
			fmt.Println("Error creating product")
		}
	}

	return nil
}
