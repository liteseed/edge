package schema

import (
	"database/sql/driver"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Status string
type StoreIds []string

const (
	Queued = "queued"
	Error  = "error"

	Sent    = "sent"
	Success = "success"
)

func (s *Status) Scan(value interface{}) error {
	*s = Status(value.(string))
	return nil
}

func (s Status) Value() (driver.Value, error) {
	return string(s), nil
}

type Order struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Status   Status    `gorm:"index:idx_status;default:queued" sql:"type:status" json:"status"`
	StoreID  uuid.UUID `json:"store_id"`
	PublicID string    `json:"public_id"`
	Checksum string    `json:"checksum"`
}
