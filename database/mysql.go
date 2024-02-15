package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMysqlDatabase(DSN string) *Query {
	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Silent),
		CreateBatchSize: 200,
	})
	if err != nil {
		panic(err)
	}
	log.Println("connect mysql db success")
	return &Query{Db: db}
}
