package model

import "time"

type LedgerEntry struct {
	ID            string    `json:"id"`
	MerchantID    string    `json:"merchant_id"`
	TransactionID string    `json:"transaction_id"`
	Credit        int64     `json:"credit"`
	Debit         int64     `json:"debit"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
}
