package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is an AWS API Gateway proxy response
type Response events.APIGatewayProxyResponse

func main() {
	switch handler := os.Getenv("HANDLER"); handler {
	case "name":
		lambda.Start(NameHandler)
	case "front":
		lambda.Start(FrontHandler)
	}
}

// NameHandler handles requests for /name etc.
func NameHandler(ctx context.Context) (Response, error) {
	var outputString, outputType, handlerName string

	name := GetName()

	if os.Getenv("OUTPUT_FORMAT") == "json" {
		handlerName = "name.ToJSON()"
		outputString = name.ToJSON()
		outputType = "application/json"
	} else {
		handlerName = "name.ToHCL()"
		outputString = name.ToHCL()
		outputType = "text/plain; charset=UTF-8"
	}
	fmt.Printf("%v", outputString)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            outputString,
		Headers: map[string]string{
			"Content-Type":      outputType,
			"X-LMHD-Func-Reply": handlerName,
		},
	}
	return resp, nil

}

// FrontHandler handles requests for /front etc.
func FrontHandler(ctx context.Context) (Response, error) {
	var outputString, outputType, handlerName string

	front := GetFront()

	if os.Getenv("OUTPUT_FORMAT") == "json" {
		handlerName = "front.ToJSON()"
		outputString = front.ToJSON()
		outputType = "application/json"
	} else {
		handlerName = "front.ToHCL()"
		outputString = front.ToHCL()
		outputType = "text/plain; charset=UTF-8"
	}
	fmt.Printf("%v", outputString)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            outputString,
		Headers: map[string]string{
			"Content-Type":      outputType,
			"X-LMHD-Func-Reply": handlerName,
		},
	}
	return resp, nil

}
