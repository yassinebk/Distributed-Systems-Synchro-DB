package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

func SeedDB(db string, whoami string) error {
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

	site, err := strconv.Atoi(whoami[2:])
	if err != nil {
		log.Panicln("[-] Error converting whoami to int")
	}

	for _, p := range products {
		p.Site = site
		_, err := productSalesRepo.CreateProduct(p)
		if err != nil {
			fmt.Println("Error creating product", err)
		}
	}

	return nil
}
