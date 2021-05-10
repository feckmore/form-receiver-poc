package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/feckmore/form-receiver-poc/model"
	"github.com/feckmore/form-receiver-poc/submit/notify"
	"github.com/feckmore/form-receiver-poc/submit/response"
)

type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call in main()
func Handler(ctx context.Context, request Request) (response.Response, error) {
	log.Printf("%+v", request)

	var wrapper model.FormWrapper
	err := json.Unmarshal([]byte(request.Body), &wrapper)
	if err != nil {
		return response.WithError(http.StatusBadRequest, err)
	}

	out, err := json.Marshal(wrapper)
	if err != nil {
		return response.WithError(http.StatusBadRequest, err)
	}
	log.Println(string(out))
	err = notify.PublishToTopic(string(out))
	if err != nil {
		return response.WithError(http.StatusBadRequest, err)
	}

	return response.WithBody(http.StatusOK, string(out), nil)
}

func main() {
	log.SetFlags(log.Llongfile)
	lambda.Start(Handler)
}
