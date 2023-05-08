package shared

import (
	"synchro-db/db"
 	"time"
)

type SentMessage struct {
	Product db.Product
	Status  string `json:"string"`
	Timestamp time.Time `json:"string"`
}