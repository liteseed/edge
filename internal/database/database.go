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

func (c *Config) GetOrders(o *schema.Order) (*[]schema.Order, error) {
	orders := &[]schema.Order{}
	err := c.DB.Where(o).Limit(25).Find(&orders).Error
	return orders, err
}

func (c *Config) UpdateOrder(o *schema.Order) error {
	return c.DB.Updates(o).Error
}

func (c *Config) UpdateOrders(orders *[]schema.Order) error {
	return c.DB.Updates(orders).Error
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
