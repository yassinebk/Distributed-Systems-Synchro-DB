package shared

import (
	"encoding/json"
	"fmt"
	"log"
	"synchro-db/db"
	"time"
)

func syncDB(message []byte, dbName string) error {
	dbConnection, err := db.ConnectToDb(dbName)

	if err != nil {
		log.Panicln("[-] Error connecting to database", dbName)
	}

	productsRepo := db.NewProductSalesRepo(dbConnection)

	var receivedMessage SentMessage

	err = json.Unmarshal(message, &receivedMessage)

	if err != nil {
		log.Panicln("[-] Error while parsing data from the wire. Check it", message)
	}

	receivedMessage.Product.ID = 0

	switch receivedMessage.Status {
	case "delete":
		_, err := productsRepo.DeleteProduct(int(receivedMessage.Product.ID))

		if err != nil {
			fmt.Println("[-] Error syncing db - operation delete - row", receivedMessage.Product)
		}
	case "create":
		fmt.Println("Here inside create")
		newProduct, err := productsRepo.CreateProduct(receivedMessage.Product)

		if err != nil {
			fmt.Println("[-] Error syncing db - operation create - row", receivedMessage.Product)
		}

		fmt.Println("[+] Success syncing db - operation create - row", newProduct)
	case "update":
		updatedProduct, err := productsRepo.UpdateProduct(receivedMessage.Product)
		if err != nil {
			fmt.Println("[-] Error syncing db - operation update - row", receivedMessage.Product)
		}

		fmt.Println("[+] Success syncing db - operation update - row", updatedProduct)

	}

	return nil
}

func BoSendUpdatedData(status string, product db.Product, whoami string) {

	toSendMessage := SentMessage{
		Product:   product,
		Status:    status,
		Timestamp: time.Now(),
	}
	connection := connect()

	jsonData, err := json.Marshal(&toSendMessage)
	if err != nil {
		log.Panicln("[-] Error marshelling data - products ")
	}

	go send(connection, fmt.Sprintf("%s-to-ho", whoami), jsonData)
}

func RecvDataFromTheWire(whoami string, updateUi func()) {

	connection := connect()

	queueName := fmt.Sprintf("%s-to-ho", whoami)
	dbName := "ho.sqlite"

	fmt.Println(dbName)

	go recv(connection, queueName, func(message []byte) {

		err := syncDB(message, dbName)
		if err != nil {
			log.Panicln("Error updating db", err)
		}
		updateUi()
		if err != nil {
			log.Panicln("[-] Error syncing database")
		}

	})

}

func SendProductsToHO(products []db.Product, whoami string) {

	connection := connect()

	queueName := fmt.Sprintf("%s-to-ho", whoami)

	for _, product := range products {
		toSendMessage := SentMessage{
			Product:   product,
			Status:    "create",
			Timestamp: time.Now(),
		}

		jsonData, err := json.Marshal(&toSendMessage)
		if err != nil {
			log.Panicln("[-] Error marshelling data - products ")
		}

		go send(connection, queueName, jsonData)
	}

}
