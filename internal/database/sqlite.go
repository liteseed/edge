package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Sqlite(url string) *Context {
	db, err := gorm.Open(sqlite.Open(url), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Silent),
		CreateBatchSize: 200,
	})
	if err != nil {
		log.Fatalln("error: database connection failed", err)
	}
	log.Println("url: " + url)
	return &Context{DB: db}

}
