package db

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey;autoincrement" json:"id"`
	ExternalID  uint           `json:"externalId"`
	Site        int            `json:"site"`
	Date        time.Time      `json:"date"`
	Product     string         `json:"product"`
	Region      string         `json:"region"`
	Qty         uint32         `json:"qty"`
	Cost        float32        `json:"cost"`
	Tax         float32        `json:"tax"`
	Sent        bool           `json:"sent" gorm:"default:false"`
	AckReceived bool           `json:"ackreceived" gorm:"default:false"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}
