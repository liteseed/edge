package database

import (
	"github.com/liteseed/edge/internal/database/schema"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Context struct {
	DB *gorm.DB
}

func New(url string) (*Context, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		CreateBatchSize: 200,
	})
	if err != nil {
		return nil, err
	}
	c := &Context{DB: db}
	err = c.Migrate()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Context) Migrate() error {
	err := c.DB.AutoMigrate(&schema.Order{})
	return err
}

func (c *Context) CreateOrder(o *schema.Order) error {
	return c.DB.Create(&o).Error
}

func (c *Context) GetOrdersByStatus(status schema.Status) (*[]schema.Order, error) {
	o := &[]schema.Order{}
	err := c.DB.Where("status = ?", status).Limit(25).Find(&o).Error
	return o, err
}

func (c *Context) UpdateStatus(id string, status schema.Status) error {
	return c.DB.Model(&schema.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (c *Context) UpdateOrder(id string, order *schema.Order) error {
	return c.DB.Model(&schema.Order{}).Where("id = ?", id).Updates(order).Error
}

func (c *Context) DeleteOrder(id string) error {
	o := &schema.Order{ID: id}
	return c.DB.Delete(o).Error
}
