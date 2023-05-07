package shared

import (
	"encoding/json"
	"fmt"
	"log"
	"synchro-db/db"
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

	switch receivedMessage.status {
	case "delete":
		_, err := productsRepo.DeleteProduct(int(receivedMessage.product.ID))

		if err != nil {
			fmt.Println("[-] Error syncing db - operation delete - row", receivedMessage.product)
		}
		break
	case "create":
		newProduct, err := productsRepo.CreateProduct(receivedMessage.product)

		if err != nil {
			fmt.Println("[-] Error syncing db - operation create - row", receivedMessage.product)
		}

		fmt.Println("[+] Success syncing db - operation create - row", newProduct)
		break
	case "update":
		updatedProduct, err := productsRepo.UpdateProduct(receivedMessage.product)
		if err != nil {
			fmt.Println("[-] Error syncing db - operation update - row", receivedMessage.product)
		}

		fmt.Println("[+] Success syncing db - operation update - row", updatedProduct)

	}

	return nil
}

func BoSendUpdatedData(status string, product db.Product, whoami string) {

	toSendMessage := SentMessage{
		product,
		status,
	}
	connection := connect()

	jsonData, err := json.Marshal(&toSendMessage)
	if err != nil {
		log.Panicln("[-] Error marshelling data - products ")
	}

	go send(connection, "bo-to-ho", jsonData)
}

func RecvDataFromTheWire(whoami string) {

	connection := connect()

	queueName := fmt.Sprintf("ho-to-%s", whoami)
	dbName := fmt.Sprintf("%s.sqlite", whoami)

	go recv(connection, queueName, func(message []byte) {

		err := syncDB(message, dbName)
		if err != nil {
			log.Panicln("[-] Error syncing database")
		}

	})

}
