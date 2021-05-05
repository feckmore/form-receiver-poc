package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request Request) (Response, error) {
	log.Printf("%+v", request)

	var wrapper FormWrapper
	err := json.Unmarshal([]byte(request.Body), &wrapper)
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	out, err := json.Marshal(wrapper)
	log.Println(string(out))

	// var buf bytes.Buffer

	// body, err := json.Marshal(map[string]interface{}{
	// 	"message": "Go Serverless v1.0! Your function executed successfully!",
	// })
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	// json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		// Body:            buf.String(),
		Body: string(out),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "form-handler",
		},
	}

	return resp, nil
}

func main() {
	log.SetFlags(log.Llongfile)
	lambda.Start(Handler)
}

type A struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type B struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
}

type FormWrapper struct {
	*A `json:",omitempty"`
	*B `json:",omitempty"`
}

func (w *FormWrapper) UnmarshalJSON(data []byte) error {
	log.Println("UnmarshalJSON()")
	log.Println(data)

	type Source struct {
		SourceID string `json:"source"`
	}

	var source Source
	err := json.Unmarshal(data, &source)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("%+v", source)

	switch source.SourceID {
	case "first":
		log.Println("do first thing")
		var a A
		return json.Unmarshal(data, &a)
	case "second":
		log.Println("do second thing")
		var b B
		err = json.Unmarshal(data, &b)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Printf("%+v", b)
		w.B = &b

		topicARN := os.Getenv("SNS_TOPIC_ARN")
		log.Println(topicARN)
		err = publishToTopic(topicARN, b.Firstname)

		return err
	default:
		log.Println("do default thing")
	}

	return nil
}

func publishToTopic(arn, message string) error {
	log.Println("publishToTopic():", arn, message)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	log.Printf("%+v", sess)
	svc := sns.New(sess)
	log.Printf("%+v", svc)
	subject := "ATS"
	result, err := svc.Publish(&sns.PublishInput{
		Subject:  &subject,
		Message:  &message,
		TopicArn: &arn,
	})

	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("%+v", *result)

	return nil
}
