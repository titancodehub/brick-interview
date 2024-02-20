package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"log"
)

type HandlerOutput struct {
	Ack bool
}

type ConsumerManager struct {
	client *sqs.Client
	url    string
}

func NewSQSConsumerManager(url string) *ConsumerManager {
	return &ConsumerManager{
		url: url,
	}
}

func (p *ConsumerManager) Init() error {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-southeast-1"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: p.url}, nil
			}),
		),
	)
	if err != nil {
		return err
	}
	p.client = sqs.NewFromConfig(cfg)
	return nil
}

func (p *ConsumerManager) Handle(ctx context.Context, handler func(message types.Message) HandlerOutput) error {
	for {
		receiveOutput, err := p.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:        &p.url,
			WaitTimeSeconds: 10,
		})

		if err != nil {
			log.Print("failed to subscribe to SQS")
			return err
		}

		for _, message := range receiveOutput.Messages {
			output := handler(message)
			if output.Ack {
				if _, err := p.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      &p.url,
					ReceiptHandle: message.ReceiptHandle,
				}); err != nil {
					log.Printf("failed to delete message from queue")
				} else {
					log.Printf("message processed")
				}
			}
		}
	}
}
