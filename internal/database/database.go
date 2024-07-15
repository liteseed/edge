package database

import (
	"errors"
	"github.com/liteseed/edge/internal/database/schema"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func New(database string, url string) (*Database, error) {
	config := &gorm.Config{CreateBatchSize: 200, Logger: logger.Default.LogMode(logger.Silent)}
	switch database {
	case "postgres":
		return Postgres(url, config)
	case "sqlite":
		return Sqlite(url, config)
	default:
		return nil, errors.New("database not supported")
	}
}

func Postgres(url string, config *gorm.Config) (*Database, error) {
	psql, err := gorm.Open(postgres.Open(url), config)
	if err != nil {
		return nil, err
	}
	db := &Database{DB: psql}
	err = db.Migrate()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Sqlite(url string, config *gorm.Config) (*Database, error) {
	database, err := gorm.Open(sqlite.Open(url), config)
	if err != nil {
		return nil, err
	}
	db := &Database{DB: database}
	err = db.Migrate()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func FromDialector(d gorm.Dialector) (*Database, error) {
	db, err := gorm.Open(d, &gorm.Config{CreateBatchSize: 200, Logger: logger.Default.LogMode(logger.Warn)})
	if err != nil {
		return nil, err
	}
	return &Database{DB: db}, nil
}

func (db *Database) Migrate() error {
	err := db.DB.AutoMigrate(&schema.Order{})
	return err
}

func (db *Database) CreateOrder(o *schema.Order) error {
	return db.DB.Create(&o).Error
}

func (db *Database) GetOrders(o *schema.Order, scopes ...Scope) (*[]schema.Order, error) {
	orders := &[]schema.Order{}
	err := db.DB.Scopes(scopes...).Where(o).Limit(10).Find(&orders).Error
	return orders, err
}

func (db *Database) GetOrder(o *schema.Order, scopes ...Scope) (*schema.Order, error) {
	order := &schema.Order{}
	err := db.DB.Scopes(scopes...).Where(o).First(&order).Error
	return order, err
}

func (db *Database) UpdateOrder(id string, o *schema.Order) error {
	return db.DB.Model(&schema.Order{}).Where("id = ?", id).Updates(o).Error
}

func (db *Database) DeleteOrder(id string) error {
	return db.DB.Delete(&schema.Order{ID: id}).Error
}

func (db *Database) Shutdown() error {
	conn, err := db.DB.DB()
	if err != nil {
		return err
	}
	return conn.Close()
}

type Scope = func(*gorm.DB) *gorm.DB

func DeadlinePassed(block int64) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("deadline > ?", block)
	}
}
