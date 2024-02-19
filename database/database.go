package database

import (
	"github.com/google/uuid"
	"github.com/liteseed/bungo/database/schema"
	"gorm.io/gorm"
)

const (
	SQLite = "sqlite"
	MySQL  = "mysql"
)

type Database struct {
	DB *gorm.DB
}

func New(path string, database string) (*Database, error) {
	db := &Database{}
	switch database {
	default:
		db = NewSqliteDatabase(path)

	}
	return db, nil
}

func (db *Database) Migrate() error {
	err := db.DB.AutoMigrate(&schema.Order{})
	return err
}

func (db *Database) CreateOrder(o *schema.Order) error {
	return db.DB.Create(&o).Error
}

func (db *Database) CreateStore(s *schema.Store) error {
	return db.DB.Create(&s).Error
}

func (db *Database) GetOrder(key string) (*schema.Order, error) {
	o := schema.Order{}
	id, err := uuid.Parse(key)
	if err != nil {
		return nil, err
	}
	err = db.DB.Where("id = ?", id).First(&o).Error
	return &o, err
}

func (db *Database) UpdateStatus(status schema.Status) error {
	return db.DB.Update("status", &status).Error
}
