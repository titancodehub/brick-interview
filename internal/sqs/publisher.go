package sqs

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type PublisherManager struct {
	client *sqs.Client
	url    string
}

func NewSQSPublisherManager(url string) *PublisherManager {
	return &PublisherManager{
		url: url,
	}
}

func (p *PublisherManager) Init() error {

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

func (p *PublisherManager) Publish(ctx context.Context, message interface{}) error {
	// parse parse json to byte
	messageByte, err := json.Marshal(message)
	if err != nil {
		return err
	}

	messageString := string(messageByte)
	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:  &messageString,
		QueueUrl:     &p.url,
		DelaySeconds: 2,
	})

	return err
}
