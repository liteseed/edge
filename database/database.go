package database

import (
	"github.com/liteseed/bungo/schema"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func New() (*Database, error) {
	return nil, nil
}

func (db *Database) Migrate() error {
	err := db.DB.AutoMigrate(&schema.Order{}, &schema.OnChainTx{}, &schema.OrderStatistic{})
	return err
}
