package notification

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func PublishToTopic(message string) error {
	topicARN := os.Getenv("SNS_TOPIC_ARN")
	log.Println("topic ARN:", topicARN)

	log.Println("publishToTopic():", topicARN, message)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	log.Printf("%+v", sess)
	svc := sns.New(sess)
	log.Printf("%+v", svc)

	result, err := svc.Publish(&sns.PublishInput{
		Subject: aws.String("Form Submission Received"), // for email subscriptions
		Message: &message,
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"Source": {
				DataType:    aws.String("String"),
				StringValue: aws.String("Source Value Goes Here"),
			},
		},
		TopicArn: &topicARN,
	})

	if err != nil {
		log.Println(err)
		return err
	}
	out, _ := json.Marshal(*result)
	log.Println(string(out))

	return nil
}
