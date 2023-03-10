package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var Db *gorm.DB

func InitDb() *gorm.DB {
	Db = connectDB("asteriskcdrdb")
	return Db
}

func connectDB(dbName string) *gorm.DB {
	Username, exists := os.LookupEnv("DB_USERNAME")
	Password, exists := os.LookupEnv("DB_PASSWORD")
	Host, exists := os.LookupEnv("DB_HOST")
	Port, exists := os.LookupEnv("DB_PORT")
	if !exists {
		fmt.Printf("NO DATA IN ENV FILE")
		return nil
	}
	var err error
	dsn := Username + ":" + Password + "@tcp" + "(" + Host + ":" + Port + ")/" + dbName + "?" + "parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("Error connecting to database : error=%v\n", err)
		return nil
	}

	return db
}
