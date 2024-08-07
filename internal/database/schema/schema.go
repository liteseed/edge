package schema

type Status string

const (
	Created = "created" // Order Created

	Queued = "queued" // Order Transaction Added
	Sent   = "sent"   // Sent to Arweave

	Confirmed = "confirmed" // > Confirmations > 25
	Release   = "release"   // Request Liteseed Reward
	Permanent = "permanent" //
	Failed    = "failed"

	Invalid = "Invalid"
)

type Order struct {
	ID             string `gorm:"primary_key;" json:"id"`
	Status         Status `gorm:"index:idx_status;default:created" sql:"type:status" json:"status"`
	TransactionID  string `json:"transaction_id"`
	BundleID       string `json:"bundle_id"`
	Size           int    `json:"size"`
	DeadlineHeight uint   `json:"deadline_height"`
}
