package schema

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

const (
	WaitOnChain    = "waiting"
	PendingOnChain = "pending"
	SuccOnChain    = "success"
	FailedOnChain  = "failed"

	// order payment status
	UnPayment      = "unpaid"
	SuccPayment    = "paid"
	ExpiredPayment = "expired"

	// ReceiptEverTx Status
	UnSpent   = "unspent"
	Spent     = "spent"
	UnRefund  = "unrefund"
	Refund    = "refunded"
	RefundErr = "refundErr"

	MaxPerOnChainSize = 2 * 1024 * 1024 * 1024 // 2 GB

	TmpFileDir = "./tmpFile"
)



type ReceiptEverTx struct {
	RawId    uint64 `grom:"primarykey"` // everTx rawId
	EverHash string `gorm:"unique"`
	Nonce    int64  // ms
	Symbol   string
	TokenTag string
	From     string
	Amount   string
	Data     string
	Sig      string

	Status string //  "unspent","spent", "unrefund", "refund"
	ErrMsg string
}

type TokenPrice struct {
	Symbol    string `gorm:"primarykey"` // token symbol
	Decimals  int
	Price     float64 // unit is USD
	ManualSet bool    // manual set
	UpdatedAt time.Time
}

type OnChainTx struct {
	gorm.Model
	ArId        string
	CurHeight   int64
	BlockId     string
	BlockHeight int64
	DataSize    string
	Reward      string         // onchain arTx reward
	Status      string         // "pending","success"
	ItemIds     datatypes.JSON // json.marshal(itemIds)
	ItemNum     int
	Kafka       bool
}
