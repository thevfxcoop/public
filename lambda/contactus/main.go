package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thevfxcoop/public/lambda/pkg/github"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Request struct {
	Text  string `json:"text"`
	Email string `json:"email"`
}

type Response struct {
	Status string `json:"status"`
}

///////////////////////////////////////////////////////////////////////////////
// HANDLER

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Route handler
	if request.Path == "/contactus" && request.HTTPMethod == http.MethodPost {
		return HandlerContactUs(ctx, request)
	}

	// By default, return NotFound
	return events.APIGatewayProxyResponse{Body: http.StatusText(http.StatusNotFound), StatusCode: http.StatusNotFound}, nil
}

func HandlerContactUs(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Request will be used to take the json response from client and build it
	var body Request
	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
	}

	// Get github client
	client := github.NewClient(ctx, request.StageVariables["GITHUB_TOKEN"], request.StageVariables["GITHUB_OWNER"])
	if client == nil {
		return events.APIGatewayProxyResponse{Body: http.StatusText(http.StatusInternalServerError), StatusCode: http.StatusInternalServerError}, nil
	}

	// List Repos
	repos, err := client.ListRepos(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadGateway}, nil
	}

	// Return response
	if response, err := json.Marshal(repos); err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
	} else {
		return events.APIGatewayProxyResponse{Body: string(response), StatusCode: http.StatusOK}, nil
	}
}

///////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	lambda.Start(Handler)
}
