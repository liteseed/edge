package database

import (
	"github.com/google/uuid"
	"github.com/liteseed/bungo/internal/database/schema"
	"gorm.io/gorm"
)

const (
	SQLite = "sqlite"
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
	err := db.Migrate()
	return db, err
}

func (db *Database) Migrate() error {
	err := db.DB.AutoMigrate(&schema.Order{})
	return err
}

func (db *Database) CreateOrder(o *schema.Order) error {
	return db.DB.Create(&o).Error
}

func (db *Database) GetOrder(id uuid.UUID) (*schema.Order, error) {
	o := &schema.Order{}
	err := db.DB.Where("id = ?", id).First(&o).Error
	return o, err
}

func (db *Database) UpdateStatus(status schema.Status) error {
	return db.DB.Update("status", &status).Error
}

func (db *Database) DeleteOrder(id uuid.UUID) error {
	o := &schema.Order{ID: id}
	return db.DB.Delete(o).Error
}
