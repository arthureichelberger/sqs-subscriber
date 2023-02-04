package sqs_subscriber

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Message struct {
	sqsClient  *sqs.Client
	Body       *string
	Attributes map[string]string

	QueueURL      string
	ReceiptHandle string
}

func (m Message) Ack(ctx context.Context) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(m.QueueURL),
		ReceiptHandle: aws.String(m.ReceiptHandle),
	}

	if _, err := m.sqsClient.DeleteMessage(ctx, input); err != nil {
		return err
	}

	return nil
}
