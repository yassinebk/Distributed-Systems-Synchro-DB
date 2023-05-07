package db

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Date      time.Time      `json:"date"`
	Product   string         `json:"product"`
	Region    string         `json:"region"`
	Qty       uint32         `json:"qty"`
	Cost      float32        `json:"cost"`
	Tax       float32        `json:"tax"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}
