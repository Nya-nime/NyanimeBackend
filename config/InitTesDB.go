package utils

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitTestDB initializes the test database
func InitTestDB() *gorm.DB {
	// Ganti dengan detail koneksi database pengujian Anda
	dsn := "root:ini_passwordnyabaruyaa14@tcp(127.0.0.1:3306)/nyanime_test?charset=utf8&parseTime=True&loc=Local"
	var err error

	// Set logger untuk debugging
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Kembalikan instance DB
	return DB
}
