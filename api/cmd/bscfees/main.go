package main

import (
	"bsc-fees/pkg/handler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return handler.HandleRequest(request)
}

func main() {
	lambda.Start(HandleRequest)
}
