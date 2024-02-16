package schema

import "time"

type Order struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	ItemId   string `gorm:"index:idx0" json:"itemId"` // bundleItem id
	Signer   string `gorm:"index:idx1" json:"signer"` // item signer
	SignType int    `json:"signType"`

	Size               int64  `json:"size"`
	Currency           string `json:"currency"` // payment token symbol
	Decimals           int    `json:"decimals"`
	Fee                string `json:"fee"`
	PaymentExpiredTime int64  `json:"paymentExpiredTime"` // uint s
	ExpectedBlock      int64  `json:"expectedBlock"`

	PaymentStatus string `gorm:"index:idx0" json:"paymentStatus"` // "unpaid", "paid", "expired"
	PaymentId     string `json:"paymentId"`                       // everHash

	Status string `gorm:"index:idx5" json:"status"` // "waiting","pending","success","failed"
}
