package service

import (
	"context"
	"github.com/titancodehub/brick-interview/client/bank"
	"github.com/titancodehub/brick-interview/model"
)

type BankService struct {
	Client bank.Client
}

func NewBankService(client bank.Client) *BankService {
	return &BankService{
		Client: client,
	}
}

func (s *BankService) ValidateAccountNumber(ctx context.Context, accountNumber string, bankCode string) (model.BankAccount, error) {
	res, err := s.Client.ValidateBankAccount(ctx, accountNumber)
	if err != nil {
		return model.BankAccount{}, err
	}

	res.BankCode = bankCode
	return res, nil
}
