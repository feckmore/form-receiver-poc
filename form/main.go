package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/feckmore/form-receiver-poc/response"
)

type FormWrapper struct {
	*A `json:",omitempty"`
	*B `json:",omitempty"`
}

type A struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type B struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
}

type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request Request) (response.Response, error) {
	log.Printf("%+v", request)

	var wrapper FormWrapper
	err := json.Unmarshal([]byte(request.Body), &wrapper)
	if err != nil {
		return response.WithError(http.StatusBadRequest, err)
	}

	out, err := json.Marshal(wrapper)
	if err != nil {
		return response.WithError(http.StatusBadRequest, err)
	}
	log.Println(string(out))
	err = publishToTopic(string(out))
	if err != nil {
		return response.WithError(http.StatusBadRequest, err)
	}

	return response.WithBody(http.StatusOK, string(out), nil)
}

func main() {
	log.SetFlags(log.Llongfile)
	lambda.Start(Handler)
}

func (w *FormWrapper) UnmarshalJSON(data []byte) error {
	log.Println("UnmarshalJSON()")
	log.Println(string(data))

	type Source struct {
		SourceID string `json:"source"`
	}

	var source Source
	err := json.Unmarshal(data, &source)
	if err != nil {
		log.Println(err)
		return err
	}
	out, _ := json.Marshal(source)
	log.Println(string(out))

	switch source.SourceID {
	case "first":
		log.Println("do first thing")
		var a A
		w.A = &a
		return json.Unmarshal(data, &a)
	case "second":
		log.Println("do second thing")
		var b B
		w.B = &b
		return json.Unmarshal(data, &b)
	default:
		log.Println("do default thing")
	}

	return nil
}

func publishToTopic(message string) error {
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
