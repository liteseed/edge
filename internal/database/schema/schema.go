package schema

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type Status string

const (
	// Order
	Created = "created" // Order Created

	Queued    = "queued"    // Order Transaction Added
	Sent      = "sent"      // Sent to Arweave
	Release   = "release"   // Request Liteseed Reward
	Permanent = "permanent" //
	Failed    = "failed"

	Invalid = "Invalid"
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
	ID             string `gorm:"primary_key;" json:"id"`
	Status         Status `gorm:"index:idx_status;default:created" sql:"type:status" json:"status"`
	TransactionID  string `json:"transaction_id"`
	BundleID       string `json:"bundle_id"`
	Size           int    `json:"size"`
	Confirmations  uint   `json:"confirmations"`
	DeadlineHeight uint   `json:"deadline_height"`
}
