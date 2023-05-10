package shared

import (
	"synchro-db/db"
	"time"
)

type SentMessage struct {
	Product   db.Product
	Site      int       `json:"site"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"time"`
}
