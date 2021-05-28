package main

import (
	"bytes"
	"encoding/json"
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

	if strings.HasSuffix(req.Path, ".yaml") {
		format = "yaml"
	}

	switch {
	case strings.HasPrefix(req.Path, "/v1/name"):
		return NameHandler(req, format)

	case strings.HasPrefix(req.Path, "/v1/front"):
		return FrontHandler(req, format)
	}

	resp := Response{
		StatusCode:      404,
		IsBase64Encoded: false,
		Body:            "NO HANDLER",
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"X-LMHD-Func-Reply":           "Handler",
			"X-LMHD-Req-String":           RequestToJSON(req),
		},
	}
	return resp, fmt.Errorf("NO HANDLER")

}

// NameHandler handles requests for /name etc.
func NameHandler(req events.APIGatewayProxyRequest, format string) (Response, error) {
	var outputString, outputType, handlerName string

	name := GetName()

	switch format {
	case "json":
		handlerName = "name.ToJSON()"
		outputString = name.ToJSON()
		outputType = "application/json"
	case "yaml":
		handlerName = "name.ToYAML()"
		outputString = name.ToYAML()
		outputType = "text/plain; charset=UTF-8"
	default:
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
			"Access-Control-Allow-Origin": "*",
			"Content-Type":                outputType,
			"X-LMHD-Func-Reply":           handlerName,
			"X-LMHD-Req-String":           RequestToJSON(req),
		},
	}
	return resp, nil

}

// FrontHandler handles requests for /front etc.
func FrontHandler(req events.APIGatewayProxyRequest, format string) (Response, error) {
	var outputString, outputType, handlerName string

	front := GetFront()
	switch format {
	case "json":
		handlerName = "front.ToJSON()"
		outputString = front.ToJSON()
		outputType = "application/json"
	case "yaml":
		handlerName = "front.ToYAML()"
		outputString = front.ToYAML()
		outputType = "text/plain; charset=UTF-8"
	default:
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
			"Access-Control-Allow-Origin": "*",
			"Content-Type":                outputType,
			"X-LMHD-Func-Reply":           handlerName,
			"X-LMHD-Req-String":           RequestToJSON(req),
		},
	}
	return resp, nil

}

// RequestToJSON outputs an API Gateway Proxy Request in JSON format
func RequestToJSON(req events.APIGatewayProxyRequest) string {
	var buf bytes.Buffer

	body, _ := json.Marshal(req)
	json.HTMLEscape(&buf, body)

	return buf.String()
}
