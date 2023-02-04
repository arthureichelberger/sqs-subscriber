package main

import (
	"context"
	"log"
	"os"
	"time"

	sqsSubscriber "github.com/arthureichelberger/sqs-subscriber"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(env("AWS_REGION", "us-east-1")))
	if err != nil {
		return
	}
	sqsClient := sqs.NewFromConfig(cfg)
	sqsSubscriberClient := sqsSubscriber.New(sqsClient)

	messageHandler := func(msg sqsSubscriber.Message) {
		defer msg.Ack(ctx)
		log.Println("Received message:", *msg.Body)

		// Do something with the message...
		return
	}

	if err := sqsSubscriberClient.Subscribe(ctx, env("SQS_QUEUE_URL", "MyQueue"), time.Second, messageHandler); err != nil {
		return
	}
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
