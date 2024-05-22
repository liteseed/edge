package schema

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type Status string

type Payment string

const (
	Queued = "queued"
	Error  = "error"

	Sent      = "sent"
	Permanent = "permanent"

	Reward = "reward"
	Done   = "done"

	Paid = "paid"
)

func (s *Status) Scan(value any) error {
	*s = Status(value.(string))
	return nil
}

func (s Status) Value() (driver.Value, error) {
	return string(s), nil
}

type Order struct {
	gorm.Model
	ID            string  `gorm:"primary_key;" json:"id"`
	Status        Status  `gorm:"index:idx_status;default:queued" sql:"type:status" json:"status"`
	TransactionId string  `json:"transaction_id"`
	Price         uint64  `json:"price"`
	Payment       Payment `gorm:"index:idx_status;default:unpaid" sql:"type:payment" json:"payment"`
}
