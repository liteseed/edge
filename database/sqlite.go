package database

import (
	"log"
	"os"
	"path"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	sqliteName = "bungo.sqlite"
)

func NewSqliteDatabase(directory string) *Query {
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open(path.Join(directory, sqliteName)), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Silent),
		CreateBatchSize: 200,
	})
	if err != nil {
		panic(err)
	}
	log.Println("connect sqlite db success")
	return &Query{Db: db}

}
