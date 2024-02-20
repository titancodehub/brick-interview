package model

import (
	"time"
)

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "PENDING"
	TransactionStatusSuccess TransactionStatus = "SUCCESS"
	TransactionStatusFailed  TransactionStatus = "FAILED"
)

type TransactionType string

const (
	TransactionTypeDisbursement TransactionType = "DISBURSEMENT"
)

type Transaction struct {
	ID            string            `json:"id"`
	MerchantID    string            `json:"merchant_id"`
	Reference     string            `json:"reference"`
	BankReference *string           `json:"-"`
	Amount        int64             `json:"amount"`
	Status        TransactionStatus `json:"status"`
	Type          TransactionType   `json:"type"`
	Metadata      JSONMap           `json:"metadata"`
	Created       time.Time         `json:"created"`
	Updated       time.Time         `json:"Updated"`
}

type TxnRequest struct {
	MerchantID    string `json:"merchant_id"`
	Amount        int64  `json:"amount"`
	Reference     string `json:"reference"`
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
}
