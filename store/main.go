package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/feckmore/form-receiver-poc/model"
)

var db *dynamodb.DynamoDB
var region, stage, table string

// Handler is our lambda handler invoked by the `lambda.Start` function call in main()
func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Println("Storage Handler():", region, stage, table)
	log.Printf("%+v", sqsEvent)

	// loop through all messages found in SQS queue
	for _, message := range sqsEvent.Records {
		log.Printf("%+v", message)

		// loop through attributes of message
		for attribute, value := range message.Attributes {
			log.Printf("%s: %s", attribute, value)
		}

		var record model.FormRecord
		err := json.Unmarshal([]byte(message.Body), &record)
		if err != nil {
			log.Println("Error unmarshalling message into record", err)
			return err
		}

		// TODO: validation
		av, err := dynamodbattribute.MarshalMap(record)
		if err != nil {
			log.Println("Error marshalling form into dynamodb attribute:", err)
			return err
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(table),
		}

		_, err = db.PutItem(input)
		if err != nil {
			log.Println("Error putting item into DyanmoDB:", err)
			return err
		}

		return nil
	}

	return nil
}

func main() {
	log.Println("main()")
	log.SetFlags(log.Llongfile)

	region = os.Getenv("REGION")
	stage = os.Getenv("STAGE")
	table = os.Getenv("DYNAMODB_TABLE_NAME")

	session, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Fatal("Failed to connect to AWS:", err)
	} else {
		db = dynamodb.New(session)
	}

	lambda.Start(Handler)
}
