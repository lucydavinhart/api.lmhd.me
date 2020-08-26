package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

func main() {
	// TODO: Other handlers, based on env var
	lambda.Start(NameHandler)
}

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
