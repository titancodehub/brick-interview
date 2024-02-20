package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/titancodehub/brick-interview/client/bank"
	"github.com/titancodehub/brick-interview/common"
	"github.com/titancodehub/brick-interview/internal/sqs"
	"github.com/titancodehub/brick-interview/model"
	"github.com/titancodehub/brick-interview/repository"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

func GetCompanyBankAccountNumber(bankCode string) string {
	switch bankCode {
	case "BNI":
		return "61110"
	case "BRI":
		return "61111"
	case "MANDIRI":
		return "61112"
	default:
		return "61113"
	}
}

type TransactionService struct {
	merchantRepository    *repository.MerchantRepository
	transactionRepository *repository.TransactionRepository
	entriesRepository     *repository.EntriesRepository
	bankClient            *bank.Client
	publisher             *sqs.PublisherManager
}

func NewTransactionService(merchantRepository *repository.MerchantRepository, transactionRepository *repository.TransactionRepository, entryRepository *repository.EntriesRepository, publisher *sqs.PublisherManager, bankClient *bank.Client) *TransactionService {
	return &TransactionService{
		merchantRepository:    merchantRepository,
		transactionRepository: transactionRepository,
		entriesRepository:     entryRepository,
		publisher:             publisher,
		bankClient:            bankClient,
	}
}

func (s *TransactionService) GetMerchantByID(ctx context.Context, id string) (model.Merchant, error) {
	return s.merchantRepository.GetByID(ctx, id)
}

func (s *TransactionService) Disbursement(ctx context.Context, req model.TxnRequest) (model.Transaction, error) {
	// This constraint should be in the database level
	_, err := s.transactionRepository.GetByReference(ctx, req.MerchantID, req.Reference)
	if err == nil {
		return model.Transaction{}, common.ErrorDuplicateReference
	}

	if !errors.Is(err, common.ErrorRecordNotExist) {
		return model.Transaction{}, common.ErrorDuplicateReference
	}

	companyBankAccount := GetCompanyBankAccountNumber(req.BankCode)
	merchantBankAccount := req.AccountNumber

	g := new(errgroup.Group)
	var merchantBankDetail model.BankAccount
	var companyBankDetail model.BankAccount
	g.Go(func() error {
		detail, goErr1 := s.bankClient.ValidateBankAccount(ctx, companyBankAccount)
		if goErr1 != nil {
			return goErr1
		}
		companyBankDetail = detail
		return nil
	})
	g.Go(func() error {
		detail, goErr2 := s.bankClient.ValidateBankAccount(ctx, merchantBankAccount)
		if goErr2 != nil {
			return goErr2
		}
		merchantBankDetail = detail
		return nil
	})

	if err = g.Wait(); err != nil {
		return model.Transaction{}, err
	}

	tx := s.merchantRepository.GetDB().Begin(&sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelSerializable,
	})

	if err := tx.Error; err != nil {
		log.Printf("failed to start transaction %v", err)
		return model.Transaction{}, err
	}

	m, err := s.merchantRepository.GetForUpdate(ctx, req.MerchantID, tx)
	if err != nil {
		if rollBackErr := tx.Rollback().Error; rollBackErr != nil {
			log.Printf("failed to rollback transaction %v", rollBackErr)
		}

		if errors.Is(err, common.ErrorRecordNotExist) {
			return model.Transaction{}, common.ErrorMerchantNotExist
		}
		return model.Transaction{}, err
	}

	newBalance := m.Balance - req.Amount
	if newBalance < 0 {
		if rollBackErr := tx.Rollback().Error; rollBackErr != nil {
			log.Printf("failed to rollback transaction %v", rollBackErr)
		}
		return model.Transaction{}, common.ErrorInsufficientBalance
	}

	err = s.merchantRepository.UpdateBalance(ctx, req.MerchantID, newBalance, tx)
	if err != nil {
		if rollBackErr := tx.Rollback().Error; rollBackErr != nil {
			log.Printf("failed to rollback transaction %v", rollBackErr)
		}
		return model.Transaction{}, err
	}

	timestamp := time.Now()

	newTransaction := model.Transaction{
		ID:         uuid.New().String(),
		MerchantID: req.MerchantID,
		Reference:  req.Reference,
		Amount:     req.Amount,
		Status:     model.TransactionStatusPending,
		Type:       model.TransactionTypeDisbursement,
		Metadata: map[string]interface{}{
			"merchant_account_number": req.AccountNumber,
			"merchant_bank_code":      req.BankCode,
			"merchant_account_name":   merchantBankDetail.AccountName,
			"company_account_number":  companyBankDetail.AccountNumber,
			"company_account_name":    companyBankDetail.AccountName,
		},
		Created: timestamp,
		Updated: timestamp,
	}
	err = s.transactionRepository.Create(ctx, &newTransaction, tx)
	if err != nil {
		if rollBackErr := tx.Rollback().Error; rollBackErr != nil {
			log.Printf("failed to rollback transaction %v", rollBackErr)
		}
		return model.Transaction{}, err
	}

	bankTransfer, err := s.bankClient.Transfer(ctx, companyBankAccount, req.AccountNumber, req.Amount)
	if err != nil {
		if rollBackErr := tx.Rollback().Error; rollBackErr != nil {
			log.Printf("failed to rollback transaction %v", rollBackErr)
		}
		return model.Transaction{}, err
	}
	newTransaction.BankReference = &bankTransfer.TransactionID
	err = s.transactionRepository.Update(ctx, &newTransaction, tx)
	if err != nil {
		if rollBackErr := tx.Rollback().Error; rollBackErr != nil {
			log.Printf("failed to rollback transaction %v", rollBackErr)
		}
		return model.Transaction{}, err
	}

	//err = s.publisher.Publish(ctx, newTransaction)
	//if err != nil {
	//	if rollBackErr := tx.Rollback().Error; rollBackErr != nil {
	//		log.Printf("failed to rollback transaction %v", rollBackErr)
	//	}
	//	return model.Transaction{}, err
	//}

	if err := tx.Commit().Error; err != nil {
		log.Printf("failed to commit transaction %v", err)
		return model.Transaction{}, err
	}

	return newTransaction, nil
}

func (s *TransactionService) DisbursementCallback(ctx context.Context, reqData model.DisbursementWebhookReq) (model.DisbursementWebhookRes, error) {
	// Ideally we always need to check signature/token from the request
	// But in this example, we don't do that for simplicity

	txn, err := s.transactionRepository.GetByBankReference(ctx, reqData.TransactionID)
	if err != nil {
		return model.DisbursementWebhookRes{}, err
	}
	txn.Status = reqData.Status
	if err = s.publisher.Publish(ctx, txn); err != nil {
		return model.DisbursementWebhookRes{}, err
	}
	return model.DisbursementWebhookRes{
		Message: "notification received successfully",
	}, nil
}

func (s *TransactionService) CompleteDisbursement(ctx context.Context, txnEvent model.Transaction) error {
	tx := s.merchantRepository.GetDB().Begin(&sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelSerializable,
	})

	if err := tx.Error; err != nil {
		log.Printf("failed to start transaction %v", err)
		return err
	}

	t, err := s.transactionRepository.GetForUpdate(ctx, txnEvent.ID, tx)
	if err != nil {
		if rollBackErr := tx.Rollback().Error; rollBackErr != nil {
			log.Printf("failed to rollback transaction %v", rollBackErr)
		}

		if errors.Is(err, common.ErrorRecordNotExist) {
			return common.ErrorMerchantNotExist
		}

		return err
	}

	if t.Status != model.TransactionStatusPending {
		if rollBackErr := tx.Rollback().Error; rollBackErr != nil {
			log.Printf("failed to rollback transaction %v", rollBackErr)
		}

		log.Printf("skip not pending transaction")
		return nil
	}

	if txnEvent.Status == model.TransactionStatusFailed {
		m, err := s.merchantRepository.GetForUpdate(ctx, t.MerchantID, tx)
		if err != nil {
			if rollBackErr := tx.Rollback().Error; rollBackErr != nil {
				log.Printf("failed to rollback transaction %v", rollBackErr)
			}
			return err
		}
		t.Status = model.TransactionStatusFailed
		// restore balance
		newBalance := m.Balance + t.Amount
		if err = s.merchantRepository.UpdateBalance(ctx, m.ID, newBalance, tx); err != nil {
			if rollBackErr := tx.Rollback().Error; rollBackErr != nil {
				log.Printf("failed to rollback transaction %v", rollBackErr)
			}
			return err
		}
	} else {
		// If status are success
		timestamp := time.Now()
		newEntry := model.LedgerEntry{
			ID:            uuid.New().String(),
			MerchantID:    t.MerchantID,
			TransactionID: t.ID,
			Credit:        t.Amount,
			Debit:         0,
			Created:       timestamp,
			Updated:       timestamp,
		}

		err = s.entriesRepository.Create(ctx, &newEntry, tx)
		if err != nil {
			if rollBackErr := tx.Rollback().Error; rollBackErr != nil {
				log.Printf("failed to rollback transaction %v", rollBackErr)
			}
			return err
		}

		t.Status = model.TransactionStatusSuccess
	}

	err = s.transactionRepository.Update(ctx, &t, tx)
	if err != nil {
		if rollBackErr := tx.Rollback().Error; rollBackErr != nil {
			log.Printf("failed to rollback transaction %v", rollBackErr)
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Print(fmt.Sprintf("failed to commit transaction %v", err))
		return err
	}

	return nil
}
