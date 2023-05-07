package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectToDb(dbName string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Product{})

	return db, nil
}
