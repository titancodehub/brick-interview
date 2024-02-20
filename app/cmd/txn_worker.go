package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/spf13/cobra"
	"github.com/titancodehub/brick-interview/config"
	"github.com/titancodehub/brick-interview/internal/postgres"
	"github.com/titancodehub/brick-interview/internal/sqs"
	"github.com/titancodehub/brick-interview/model"
	"github.com/titancodehub/brick-interview/repository"
	"github.com/titancodehub/brick-interview/service"
	"log"
)

var StartWorker = &cobra.Command{
	Use:   "start-worker",
	Short: "Start worker",
	Long:  "Start worker",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := postgres.CreateConnection()
		if err != nil {
			panic(fmt.Sprintf("failed to connect to the database %v", err))
		}

		// Init Repository
		merchantRepo := repository.NewMerchantRepository(db)
		transactionRepo := repository.NewTransactionRepository(db)
		entriesRepo := repository.NewEntriesRepository(db)

		// Init Service
		transactionService := service.NewTransactionService(merchantRepo, transactionRepo, entriesRepo, nil, nil)

		consumer := sqs.NewSQSConsumerManager(config.GetSQSUrl())
		err = consumer.Init()
		if err != nil {
			panic("failed to init consumer")
		}
		ctx := context.Background()

		err = consumer.Handle(ctx, func(message types.Message) sqs.HandlerOutput {
			var txn model.Transaction
			err = json.Unmarshal([]byte(*message.Body), &txn)
			if err != nil {
				log.Printf("failed to process message %v , error: %v", message, err)
				return sqs.HandlerOutput{}
			}

			log.Printf("handling %v", txn)
			err := transactionService.CompleteDisbursement(ctx, txn)
			if err != nil {
				// We can implement fallback scenario here
				// Example: exponential backoff, and routing do dlq
				log.Printf("failed handling message %v", err)
			}

			return sqs.HandlerOutput{
				Ack: true,
			}
		})

		if err != nil {
			panic(fmt.Sprintf("failed to start consumer %v", err))
		}
	},
}
