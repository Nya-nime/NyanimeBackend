package tes

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitTestDB() {
	var err error
	dsn := "root:ini_passwordnyabaruyaa14@tcp(127.0.0.1:3306)/nyanime_test?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

}
