package models

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDsn() (string, *gorm.DB) {
	dbname := os.Getenv("DATABASE_NAME")

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DATABASE_OWNER"),
		os.Getenv("DATABASE_PASSWORD"),
		dbname,
		os.Getenv("DATABASE_PORT"))

	// get connect db variable
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("連接資料庫失敗: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Fatalf("無法取得 sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	return dbname, db
}
