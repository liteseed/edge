package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Postgres(url string) *Context {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Silent),
		CreateBatchSize: 200,
	})
	if err != nil {
		log.Fatalln("error: database connection failed", err)
	}
	log.Println("database: postgres: ", url)
	return &Context{DB: db}
}
