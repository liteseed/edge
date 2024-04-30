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

func (c *Config) GetOrdersByStatus(status schema.Status) (*[]schema.Order, error) {
	o := &[]schema.Order{}
	err := c.DB.Where("status = ?", status).Limit(25).Find(&o).Error
	return o, err
}

func (c *Config) UpdateStatus(id string, status schema.Status) error {
	return c.DB.Model(&schema.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (c *Config) UpdateOrder(id string, order *schema.Order) error {
	return c.DB.Model(&schema.Order{}).Where("id = ?", id).Updates(order).Error
}

func (c *Config) DeleteOrder(id string) error {
	o := &schema.Order{ID: id}
	return c.DB.Delete(o).Error
}

func (c *Config) Shutdown() error {
	db, err := c.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
