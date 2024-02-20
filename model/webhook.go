package model

type DisbursementWebhookReq struct {
	TransactionID string            `json:"transaction_id"`
	Status        TransactionStatus `json:"status"`
}

type DisbursementWebhookRes struct {
	Message string `json:"message"`
}
