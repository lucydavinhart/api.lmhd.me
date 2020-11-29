package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is an AWS API Gateway proxy response
type Response events.APIGatewayProxyResponse

func main() {
	lambda.Start(Handler)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(req events.APIGatewayProxyRequest) (Response, error) {
	format := "hcl"

	if strings.HasSuffix(req.Path, ".json") {
		format = "json"
	}

	switch {
	case strings.HasPrefix(req.Path, "/v1/name"):
		return NameHandler(req, format)

	case strings.HasPrefix(req.Path, "/v1/front"):
		return FrontHandler(req, format)
	}

	return Response{}, fmt.Errorf("NO HANDLER")

}

// NameHandler handles requests for /name etc.
func NameHandler(req events.APIGatewayProxyRequest, format string) (Response, error) {
	var outputString, outputType, handlerName string

	name := GetName()

	if format == "json" {
		handlerName = "name.ToJSON()"
		outputString = name.ToJSON()
		outputType = "application/json"
	} else {
		handlerName = "name.ToHCL()"
		outputString = name.ToHCL()
		outputType = "text/plain; charset=UTF-8"
	}
	reqString := fmt.Sprintf("%v", req)
	fmt.Printf("%v", outputString)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            outputString,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Content-Type":                outputType,
			"X-LMHD-Func-Reply":           handlerName,
			"X-LMHD-Req-String":           reqString,
		},
	}
	return resp, nil

}

// FrontHandler handles requests for /front etc.
func FrontHandler(req events.APIGatewayProxyRequest, format string) (Response, error) {
	var outputString, outputType, handlerName string

	front := GetFront()

	if format == "json" {
		handlerName = "front.ToJSON()"
		outputString = front.ToJSON()
		outputType = "application/json"
	} else {
		handlerName = "front.ToHCL()"
		outputString = front.ToHCL()
		outputType = "text/plain; charset=UTF-8"
	}
	reqString := fmt.Sprintf("%v", req)
	fmt.Printf("%v", outputString)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            outputString,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Content-Type":                outputType,
			"X-LMHD-Func-Reply":           handlerName,
			"X-LMHD-Req-String":           reqString,
		},
	}
	return resp, nil

}
