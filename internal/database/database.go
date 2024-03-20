package database

import (
	"github.com/google/uuid"
	"github.com/liteseed/edge/internal/database/schema"
	"gorm.io/gorm"
)

const (
	SQLite = "sqlite"
)

type Database struct {
	DB *gorm.DB
}

func New(options string, path string) (*Database, error) {
	db := &Database{}
	switch options {
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

func (db *Database) GetQueuedOrders(limit int) (*[]schema.Order, error) {
	o := &[]schema.Order{}
	err := db.DB.Where("status = ?", schema.Queued).Find(&o).Limit(limit).Error
	return o, err
}

func (db *Database) UpdateStatus(id uuid.UUID, status schema.Status) error {
	return db.DB.Model(&schema.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (db *Database) DeleteOrder(id uuid.UUID) error {
	o := &schema.Order{ID: id}
	return db.DB.Delete(o).Error
}
