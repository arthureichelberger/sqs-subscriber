package sqs_subscriber

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Client struct {
	client *sqs.Client
}

func New(client *sqs.Client) *Client {
	return &Client{
		client: client,
	}
}

func (s *Client) Subscribe(ctx context.Context, queueName string, pollingFrequency time.Duration, handler func(msg Message)) error {
	queue, err := s.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: &queueName})
	if err != nil {
		return err
	}

	sink := make(chan error, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				sink <- nil
				return
			default:
				func() {
					defer time.Sleep(pollingFrequency)

					messages, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{QueueUrl: queue.QueueUrl, MaxNumberOfMessages: 10})
					if err != nil {
						sink <- err
						return
					}

					if len(messages.Messages) == 0 {
						return
					}

					for _, message := range messages.Messages {
						msg := Message{
							sqsClient:     s.client,
							Body:          message.Body,
							Attributes:    message.Attributes,
							QueueURL:      *queue.QueueUrl,
							ReceiptHandle: *message.ReceiptHandle,
						}

						go handler(msg)
					}
				}()

			}
		}
	}()

	if err := <-sink; err != nil {
		return err
	}

	return nil
}
