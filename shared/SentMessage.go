package shared

import "synchro-db/db"

type SentMessage struct {
	product db.Product
	status  string `json:"string"`
}
