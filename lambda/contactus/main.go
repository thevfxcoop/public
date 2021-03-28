package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thevfxcoop/public/lambda/pkg/github"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Request struct {
	Repository string `json:"repository"`
	Text       string `json:"text"`
	Email      string `json:"email"`
}

type Response struct {
	Status  string      `json:"status"`
	Content interface{} `json:"content,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// HANDLER

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// COORS
	if request.HTTPMethod == http.MethodOptions {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "https://www.vfx.coop",
				"Access-Control-Allow-Methods": "POST, GET, OPTIONS",
				"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
			},
			StatusCode: http.StatusOK,
		}, nil
	}

	// Route handler
	if request.Path == "/contactus" && request.HTTPMethod == http.MethodPost {
		return HandlerContactUs(ctx, request)
	}

	// By default, return NotFound
	return SendResponse(http.StatusNotFound, nil)
}

func HandlerContactUs(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Request will be used to take the json response from client and build it
	var body Request
	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		return SendResponse(http.StatusBadRequest, err.Error())
	}

	// Check request
	if body.isValid() == false {
		return SendResponse(http.StatusBadRequest, nil)
	}

	// Get github client
	client := github.NewClient(ctx, request.StageVariables["GITHUB_TOKEN"], request.StageVariables["GITHUB_OWNER"])
	if client == nil {
		return SendResponse(http.StatusInternalServerError, nil)
	}

	// Create an issue
	_, err := client.CreateIssue(ctx, body.Repository, body.Email, body.Text+"\n")
	if err != nil {
		return SendResponse(http.StatusBadGateway, err.Error())
	}

	// Return response
	return SendResponse(http.StatusOK, "Thanks for contacting us")
}

///////////////////////////////////////////////////////////////////////////////
// RESPONDER

func SendResponse(code int, value interface{}) (events.APIGatewayProxyResponse, error) {
	response := new(Response)
	if value == nil {
		response.Status = http.StatusText(code)
	} else if code == http.StatusOK {
		response.Status = http.StatusText(code)
		response.Content = value
	} else {
		response.Status = fmt.Sprint(value)
	}
	json, err := json.Marshal(response)
	return events.APIGatewayProxyResponse{
		Body:       string(json) + "\n",
		StatusCode: code,
	}, err
}

///////////////////////////////////////////////////////////////////////////////
// UTILS

func (r Request) isValid() bool {
	if r.Email == "" {
		return false
	}
	if r.Text == "" {
		return false
	}
	if r.Repository == "" {
		return false
	}
	return true
}

///////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	lambda.Start(Handler)
}
