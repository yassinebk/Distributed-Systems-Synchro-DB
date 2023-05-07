package shared

import (
	"encoding/json"
	"fmt"
	"log"
	"synchro-db/db"
)

func SendUpdatedData(products []db.Product, whoami string) {

	connection := connect()

	jsonData, err := json.Marshal(products)
	if err != nil {
		log.Panicln("[-] Error marshelling data - products ")
	}

	go send(connection, fmt.Sprintf("%s-to-ho", whoami), jsonData)
}

func RecvDataFromTheWire(whoami string) {

	connection := connect()

	queueName := fmt.Sprintf("ho-to-%s", whoami)
	dbName := fmt.Sprintf("%s.sqlite", whoami)

	go recv(connection, queueName, func(message []byte) {

		dbConnection, err := db.ConnectToDb(dbName)
		if err != nil {
			log.Panicln("[-] Error connecting to db of ", dbName)
		}

		productRepos := db.NewProductSalesRepo(dbConnection)

		var products []db.Product

		err = json.Unmarshal(message, &products)
		if err != nil {

			log.Panicln("[-] Error marshelling products ", err)
		}

		productRepos.BatchUpsert(products)

	})

}

func PingForUpdates() {

}
