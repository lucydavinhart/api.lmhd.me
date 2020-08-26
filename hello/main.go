package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rodaine/hclencoder"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type Res struct {
	Message string `hcl:"message"`
}

// HandlerJSON is our lambda handler invoked by the `lambda.Start` function call
func HandlerJSON(ctx context.Context) (Response, error) {
	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"message": "Go Serverless v1.0! Your function executed successfully!",
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

// HandlerHCL is our lambda handler invoked by the `lambda.Start` function call
func HandlerHCL(ctx context.Context) (Response, error) {

	r := Res{
		Message: "Go Serverless v1.0! Your function executed successfully!",
	}

	hcl, err := hclencoder.Encode(r)
	if err != nil {
		log.Fatal("unable to encode: ", err)
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(hcl),
		Headers: map[string]string{
			"Content-Type":      "application/json",
			"X-LMHD-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	if os.Getenv("OUTPUT_FORMAT") == "json" {
		lambda.Start(HandlerJSON)
	} else {
		lambda.Start(HandlerHCL)
	}
}
