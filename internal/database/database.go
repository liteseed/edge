package database

import (
	"github.com/liteseed/edge/internal/database/schema"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DB *gorm.DB
}

func New(url string) (*Config, error) {
	db, err := gorm.Open(sqlite.Open(url), &gorm.Config{
		CreateBatchSize: 200,
		Logger:          logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	c := &Config{DB: db}
	err = c.Migrate()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) Migrate() error {
	err := c.DB.AutoMigrate(&schema.Order{})
	return err
}

func (c *Config) CreateOrder(o *schema.Order) error {
	return c.DB.Create(&o).Error
}

func (c *Config) GetOrders(o *schema.Order, scopes ...Scope) (*[]schema.Order, error) {
	orders := &[]schema.Order{}
	err := c.DB.Scopes(scopes...).Where(o).Limit(25).Find(&orders).Error
	return orders, err
}

func (c *Config) GetOrder(o *schema.Order, scopes ...Scope) (*schema.Order, error) {
	order := &schema.Order{}
	err := c.DB.Scopes(scopes...).Where(o).First(&order).Error
	return order, err
}

func (c *Config) UpdateOrder(o *schema.Order) error {
	return c.DB.Model(schema.Order{ID: o.ID}).Updates(o).Error
}

func (c *Config) DeleteOrder(id string) error {
	return c.DB.Delete(&schema.Order{ID: id}).Error
}

func (c *Config) Shutdown() error {
	db, err := c.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

type Scope = func(*gorm.DB) *gorm.DB

// Scope for filtering records where confirmations is greater than 25
func ConfirmationsGreaterThanEqualTo25(db *gorm.DB) *gorm.DB {
	return db.Where("confirmations >= ?", 25)
}

// Scope for filtering records where confirmations is greater than 25
func ConfirmationsLessThan25(db *gorm.DB) *gorm.DB {
	return db.Where("confirmations < ?", 25)
}

func DeadlinePassed(block int64) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("deadline > ?", block)
	}
}
