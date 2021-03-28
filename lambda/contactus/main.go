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
	Status  int         `json:"status"`
	Content interface{} `json:"content,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	ResponseHeaders = map[string]string{
		"Access-Control-Allow-Origin":  "https://www.vfx.coop",
		"Access-Control-Allow-Methods": "POST, GET, OPTIONS",
		"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	}
)

///////////////////////////////////////////////////////////////////////////////
// HANDLER

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// COORS
	if request.HTTPMethod == http.MethodOptions {
		return events.APIGatewayProxyResponse{
			Headers:    ResponseHeaders,
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
	_, err := client.CreateIssue(ctx, body.Repository, "Message from "+body.Email, body.Text+"\n", []string{"contactus"})
	if err != nil {
		return SendResponse(http.StatusBadGateway, err.Error())
	}

	// Return response
	return SendResponse(http.StatusOK, "Thanks for contacting us, we'll be in touch shortly")
}

///////////////////////////////////////////////////////////////////////////////
// RESPONDER

func SendResponse(code int, value interface{}) (events.APIGatewayProxyResponse, error) {
	response := new(Response)
	response.Status = code
	if code == http.StatusOK {
		response.Content = value
	} else {
		response.Content = fmt.Sprint(value)
	}
	json, err := json.Marshal(response)
	return events.APIGatewayProxyResponse{
		Body:       string(json) + "\n",
		Headers:    ResponseHeaders,
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
