package database

import (
	"github.com/liteseed/bungo/internal/database/schema"
	"gorm.io/gorm"
)

const (
	SQLite     = "sqlite"
	MySQL      = "mysql"
	PostgreSQL = "postgres"
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
	err := db.DB.AutoMigrate(&schema.Order{}, &schema.Store{})
	return err
}

func (db *Database) CreateOrder(o *schema.Order) error {
	return db.DB.Create(&o).Error
}

func (db *Database) CreateStore(s *schema.Store) error {
	return db.DB.Create(&s).Error
}

func (db *Database) GetOrder(id string) (*schema.Order, error) {
	o := &schema.Order{}
	err := db.DB.Where("id = ?", id).First(&o).Error
	return o, err
}

func (db *Database) GetStores(oid string) (*[]schema.Store, error) {
	s := &[]schema.Store{}
	err := db.DB.Where("id = ?", oid).Find(&s).Error
	return s, err
}

func (db *Database) UpdateStatus(status schema.Status) error {
	return db.DB.Update("status", &status).Error
}
