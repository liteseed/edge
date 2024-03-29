package database

import (
	"github.com/liteseed/edge/internal/database/schema"
	"gorm.io/gorm"
)

const (
	SQLite     = "sqlite"
	PostgreSQL = "postgresql"
)

type Context struct {
	DB *gorm.DB
}

func New(options string, url string) (*Context, error) {
	var c *Context
	switch options {
	case "postgres":
		c = Postgres(url)
	default:
		c = Sqlite(url)
	}
	err := c.Migrate()
	return c, err
}

func (c *Context) Migrate() error {
	err := c.DB.AutoMigrate(&schema.Order{})
	return err
}

func (c *Context) CreateOrder(o *schema.Order) error {
	return c.DB.Create(&o).Error
}

func (c *Context) GetOrder(id string) (*schema.Order, error) {
	o := &schema.Order{}
	err := c.DB.Where("public_id = ?", id).First(&o).Error
	return o, err
}

func (c *Context) GetOrdersByStatus(status schema.Status) (*[]schema.Order, error) {
	o := &[]schema.Order{}
	err := c.DB.Where("status = ?", status).Limit(25).Find(&o).Error
	return o, err
}

func (c *Context) UpdateStatus(id string, status schema.Status) error {
	return c.DB.Model(&schema.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (c *Context) UpdateTransactionID(id string, transactionId string) error {
	return c.DB.Model(&schema.Order{}).Where("id = ?", id).Update("transaction_id", transactionId).Error
}

func (c *Context) DeleteOrder(id string) error {
	o := &schema.Order{ID: id}
	return c.DB.Delete(o).Error
}
