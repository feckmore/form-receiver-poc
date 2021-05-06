package response

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Error contains info for an error response
type Error struct {
	Err     error  `json:"-"`
	Message string `json:"message"`
	Code    int    `json:"status"`
}

type Response events.APIGatewayProxyResponse

// Error fulfills the Go error interface for the Error struct
func (err Error) Error() string {
	if err.Message != "" {
		return err.Message
	}

	return err.Err.Error()
}

func WithBody(statusCode int, body string, err error) (Response, error) {
	if err != nil {
		return WithError(statusCode, err)
	}

	resp := Response{
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Body:            body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, err
}

func WithError(statusCode int, err error) (Response, error) {
	log.Println("errorResponse()", statusCode, err)
	// check if error is our custom Error type
	rErr, ok := err.(*Error)
	if !ok {
		if e, ok := err.(Error); ok {
			rErr = &e
		}
	}
	// there's an error, but not in our customer error format
	if rErr == nil {
		if err == nil {
			err = errors.New(http.StatusText(http.StatusInternalServerError))
		}
		rErr = &Error{Err: err}
	}

	// validate status code: if not specified or invalid, use 500
	if http.StatusText(rErr.Code) == "" { // code within error is primary
		if http.StatusText(statusCode) != "" { // code passed into function is secondary
			rErr.Code = statusCode
		} else {
			rErr.Code = http.StatusInternalServerError
		}
	}

	// create body of error message
	var buf bytes.Buffer
	body, _ := json.Marshal(map[string]interface{}{
		"message": rErr.Err.Error(),
	})
	json.HTMLEscape(&buf, body)

	return WithBody(rErr.Code, string(body), nil)
}
