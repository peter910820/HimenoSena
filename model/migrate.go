package model

import (
	"gorm.io/gorm"
)

func Migration(dbName string, db *gorm.DB) {
	db.AutoMigrate(&Member{})
}
