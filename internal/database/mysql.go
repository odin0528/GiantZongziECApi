package database

import (
	"fmt"
	"log"
	"os"

	env "github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 全域
var DB *gorm.DB

func init() {
	env.Load()
	var err error
	DB, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
}
