package database

import (
	"log"

	"github.com/liteseed/edge/internal/database/schema"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Context struct {
	DB *gorm.DB
}

func New(url string) (*Context, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Silent),
		CreateBatchSize: 200,
	})
	if err != nil {
		log.Fatalln("error: database connection failed", err)
	}
	log.Println("url: ", url)
	c := &Context{DB: db}
	err = c.Migrate()
	if err != nil {
		log.Fatalln("error: database connection failed", err)
	}
	return c, err
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

func (c *Context) UpdateTransactionID(id string, transactionId string) error {
	return c.DB.Model(&schema.Order{}).Where("id = ?", id).Update("transaction_id", transactionId).Error
}

func (c *Context) DeleteOrder(id string) error {
	o := &schema.Order{ID: id}
	return c.DB.Delete(o).Error
}
