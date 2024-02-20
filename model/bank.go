package model

type BankAccount struct {
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	BankCode      string `json:"bank_code"`
}

type BankTransfer struct {
	SourceAccountNumber      string `json:"source_account_number"`
	DestinationAccountNumber string `json:"destination_account_number"`
	Amount                   int64  `json:"amount"`
	Status                   string `json:"status"`
	TransactionID            string `json:"transaction_id"`
}
