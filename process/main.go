package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call in main()
func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Printf("%+v", sqsEvent)

	// loop through all messages found in SQS queue
	for _, message := range sqsEvent.Records {
		log.Printf("%+v", message)
		// loop through attributes of message
		for attribute, value := range message.Attributes {
			log.Printf("%s: %s", attribute, value)
		}
	}

	return nil
}

func main() {
	log.Println("main()")
	log.SetFlags(log.Llongfile)
	lambda.Start(Handler)
}
