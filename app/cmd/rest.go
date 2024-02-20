package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/titancodehub/brick-interview/client/bank"
	"github.com/titancodehub/brick-interview/common"
	"github.com/titancodehub/brick-interview/config"
	"github.com/titancodehub/brick-interview/internal/postgres"
	"github.com/titancodehub/brick-interview/internal/sqs"
	"github.com/titancodehub/brick-interview/middleware"
	"github.com/titancodehub/brick-interview/model"
	"github.com/titancodehub/brick-interview/repository"
	"github.com/titancodehub/brick-interview/service"
	"net/http"
)

var StartRestServer = &cobra.Command{
	Use:   "start-rest-server",
	Short: "Start rest server",
	Long:  "Start rest server",
	Run: func(cmd *cobra.Command, args []string) {
		r :=
			gin.Default()
		r.Use(middleware.RequestValidation())

		db, err := postgres.CreateConnection()
		if err != nil {
			panic(fmt.Sprintf("failed to connect to the database %v", err))
		}

		// Init Repository
		merchantRepo := repository.NewMerchantRepository(db)
		transactionRepo := repository.NewTransactionRepository(db)
		entriesRepo := repository.NewEntriesRepository(db)

		// Init Publisher
		publisher := sqs.NewSQSPublisherManager(config.GetSQSUrl())
		err = publisher.Init()
		if err != nil {
			panic("failed to init publisher")
		}

		// Init Client
		bankClient := bank.NewClient(config.GetBankClientURL())

		// Init Service
		transactionService := service.NewTransactionService(merchantRepo, transactionRepo, entriesRepo, publisher, &bankClient)
		bankService := service.NewBankService(bankClient)

		// Init Routes
		r.GET("/validate-bank-accounts", func(c *gin.Context) {
			accountNumber := c.Query("account_number")
			bankCode := c.Query("bank_code")
			result, err := bankService.ValidateAccountNumber(c, accountNumber, bankCode)
			if err != nil {
				c.Error(err)
				return
			}

			c.JSON(http.StatusOK, result)
		}, middleware.ErrorHandler(map[error]int{
			common.ErrorBankAccountNotExist: http.StatusNotFound,
		}))

		r.POST("/disbursements", func(c *gin.Context) {
			var reqData model.TxnRequest
			err := c.Bind(&reqData)

			if err != nil {
				c.Error(err)
				return
			}

			result, err := transactionService.Disbursement(c, reqData)
			if err != nil {
				c.Error(err)
				return
			}

			c.JSON(http.StatusOK, result)
		}, middleware.ErrorHandler(map[error]int{
			common.ErrorRecordNotExist:      http.StatusNotFound,
			common.ErrorInsufficientBalance: http.StatusForbidden,
			common.ErrorMerchantNotExist:    http.StatusNotFound,
			common.ErrorDuplicateReference:  http.StatusConflict,
			common.ErrorBankAccountNotExist: http.StatusNotFound,
		}))

		r.POST("/webhooks/disbursements", func(c *gin.Context) {
			var reqData model.DisbursementWebhookReq
			err := c.Bind(&reqData)

			if err != nil {
				c.Error(err)
				return
			}

			result, err := transactionService.DisbursementCallback(c, reqData)
			if err != nil {
				c.Error(err)
				return
			}

			c.JSON(http.StatusOK, result)
		}, middleware.ErrorHandler(map[error]int{
			common.ErrorRecordNotExist:      http.StatusNotFound,
			common.ErrorInsufficientBalance: http.StatusForbidden,
			common.ErrorMerchantNotExist:    http.StatusNotFound,
			common.ErrorDuplicateReference:  http.StatusConflict,
			common.ErrorBankAccountNotExist: http.StatusNotFound,
		}))

		err = r.Run()
		if err != nil {
			panic("Error")
		}
	},
}
