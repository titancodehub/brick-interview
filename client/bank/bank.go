package bank

import (
	"context"
	"fmt"
	"github.com/monaco-io/request"
	"github.com/titancodehub/brick-interview/common"
	"github.com/titancodehub/brick-interview/model"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) Client {
	return Client{
		BaseURL: baseURL,
	}
}

func (c *Client) ValidateBankAccount(ctx context.Context, accountNumber string) (model.BankAccount, error) {

	var result model.BankAccount
	resp := request.NewWithContext(ctx).
		GET(fmt.Sprintf("%s/api/bank-accounts/%s", c.BaseURL, accountNumber)).
		Send()

	if !resp.OK() {
		return result, resp.Error()
	}

	if resp.Code() == 404 {
		return result, common.ErrorBankAccountNotExist
	}

	if resp.Code() > 199 && resp.Code() < 300 {
		resp = resp.Scan(&result)
	}

	if resp.Error() != nil {
		return result, resp.Error()
	}

	return result, nil
}

func (c *Client) Transfer(ctx context.Context, srcAccountNumber, destAccountNumber string, amount int64) (model.BankTransfer, error) {
	// Transfer API are mocked, it will only return the id, we assume that the status will always be PENDING
	var respBody struct {
		TransactionID            string `json:"transaction_id"`
		SourceAccountNumber      string `json:"source_account_number"`
		DestinationAccountNumber string `json:"destination_account_number"`
		Amount                   int64  `json:"amount"`
	}

	var result model.BankTransfer
	resp := request.NewWithContext(ctx).
		POST(fmt.Sprintf("%s/api/transfers", c.BaseURL)).
		AddJSON(map[string]interface{}{
			"source_account_number":      srcAccountNumber,
			"destination_account_number": destAccountNumber,
			"amount":                     amount,
		}).
		Send()

	if !resp.OK() {
		return result, resp.Error()
	}

	if resp.Code() > 199 && resp.Code() < 300 {
		resp = resp.Scan(&respBody)
	}

	if resp.Error() != nil {
		return result, resp.Error()
	}

	result.SourceAccountNumber = respBody.SourceAccountNumber
	result.DestinationAccountNumber = respBody.DestinationAccountNumber
	result.Amount = respBody.Amount
	result.TransactionID = respBody.TransactionID
	result.Status = "PENDING"

	return result, nil
}
