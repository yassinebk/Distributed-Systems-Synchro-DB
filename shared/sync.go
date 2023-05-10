// package shared is the package that holds all the methods and global variables
// needed to sync data over the wire without messing it up
package shared

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"synchro-db/db"
	"time"

	"github.com/fatih/color"
)

var (
	productLocalQueu         MinHeap
	unorderedMessagesChannel = make(chan SentMessage, 10)
	orderedMessagesChannel   = make(chan int, 10)
)

// syncDB is a private function that is responsible of identifying the database operation in the function and performing it
func syncDB(receivedMessage SentMessage, dbName string) error {
	dbConnection, err := db.ConnectToDb(dbName)

	if err != nil {
		log.Panicln("[-] Error connecting to database", dbName)
	}

	productsRepo := db.NewProductSalesRepo(dbConnection)

	switch receivedMessage.Status {
	case "delete":
		_, err := productsRepo.DeleteProduct(int(receivedMessage.Product.ID))

		if err != nil {
			log.Println("[-] Error syncing db - operation delete - row", receivedMessage.Product)
		}
	case "create":
		receivedMessage.Product.ID = 0
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

// BoSendUpdatedData notifies the HO of data update
func BoSendUpdatedData(status string, product db.Product, whoami string) {

	site, err := strconv.Atoi(whoami[1:])
	if err != nil {
		log.Panicln("[-] Error converting whoami to int")
	}
	log.Panicln("site", site)
	log.Panicln("whoami", whoami)

	toSendMessage := SentMessage{
		Product:   product,
		Status:    status,
		Site:      site,
		Timestamp: time.Now(),
	}
	connection := connect()

	jsonData, err := json.Marshal(&toSendMessage)
	if err != nil {
		log.Panicln("[-] Error marshelling data - products ")
	}

	go send(connection, fmt.Sprintf("%s-to-ho", whoami), jsonData)
}

func RecvDataFromTheWire(whoami string) {

	connection := connect()

	queueName := fmt.Sprintf("%s-to-ho", whoami)

	go recv(connection, queueName, func(message []byte) {
		var receivedMessage SentMessage

		err := json.Unmarshal(message, &receivedMessage)

		if err != nil {
			log.Panicln("[-] Error while parsing data from the wire. Check it", message)
		}

		unorderedMessagesChannel <- receivedMessage

	})

}

// OrderProductIntoHeap  is our second (goroutine / thread ). It received a new message from from the main thread
// and organize it inside the heap then notifies the third thread to sync it into the db
func OrderProductIntoHeap() {

	connection := connect()
	queueName := "ho-to-bo"

	magenta := color.New(color.FgMagenta).SprintFunc()
	fmt.Println(magenta("[+] Second thread started workin"))
	for {

		select {
		case receivedMessage := <-unorderedMessagesChannel:
			productLocalQueu.Push(receivedMessage)

			receivedMessage.Product.AckReceived = true

			jsonData, err := json.Marshal(&receivedMessage)
			if err != nil {
				log.Panicln("[-] Error marshelling data - products ")
			}
			go send(connection, queueName+fmt.Sprintf("%d", receivedMessage.Site), jsonData)

			orderedMessagesChannel <- 1
		default:
			fmt.Println(magenta("[~] WAITING FOR MORE PRODUCTS TO PUT IN HEAP"))
			time.Sleep(10 * time.Second)
		}
	}
}

// PerformDbOp when notified of a new message it pops the heap and performs the operation
func PerformDbOp(updateUi func()) {

	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Println(yellow("[+] Third thread working"))
	for {

		select {
		case <-orderedMessagesChannel:
			receivedMessage := productLocalQueu.Pop()
			err := syncDB(receivedMessage.(SentMessage), "ho.sqlite")
			if err != nil {
				log.Panicln("Error updating db", err)
			}
			updateUi()
			if err != nil {
				log.Panicln("[-] Error syncing database")
			}
		default:
			fmt.Println(yellow("[~] WAITING FOR MORE DB OPERATIONS"))
			time.Sleep(10 * time.Second)
		}
	}

}

func SendProductsToHO(products []db.Product, whoami string) {

	connection := connect()
	queueName := fmt.Sprintf("%s-to-ho", whoami)

	site, err := strconv.Atoi(whoami[2:])
	if err != nil {
		log.Panicln("[-] Error converting whoami to int")
	}

	for _, product := range products {
		toSendMessage := SentMessage{
			Product:   product,
			Status:    "create",
			Site:      site,
			Timestamp: time.Now(),
		}

		jsonData, err := json.Marshal(&toSendMessage)
		if err != nil {
			log.Panicln("[-] Error marshelling data - products ")
		}

		go send(connection, queueName, jsonData)
	}
}

func ListenOnAcks(whoami string, productRepo *db.ProductSalesRepo) {

	dbConnection, err := db.ConnectToDb(whoami + ".sqlite") //

	if err != nil {
		log.Panicln("[-] Error connecting to database", err)
	}

	productsRepo := db.NewProductSalesRepo(dbConnection)

	connection := connect()

	queueName := "ho-to-bo" + whoami[2:]

	go recv(connection, queueName, func(message []byte) {
		var receivedMessage SentMessage

		err := json.Unmarshal(message, &receivedMessage)

		if err != nil {
			log.Panicln("[-] Error while parsing data from the wire. Check it", message)
		}

		product := receivedMessage.Product
		product.AckReceived = true

		productsRepo.UpdateProduct(product)

	})

}
