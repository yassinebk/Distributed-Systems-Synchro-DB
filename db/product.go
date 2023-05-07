package db

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID uint `gorm:"primaryKey"`

	Date      time.Time
	Product   string
	Region    string
	Qty       uint32
	Cost      float32
	Tax       float32
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
