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

	case strings.HasPrefix(req.Path, "/v1/quest"):
		return QuestHandler(req, format)

	case strings.HasPrefix(req.Path, "/v1/auth"):
		return AuthHandler(req, format)

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

// RequestToJSON outputs an API Gateway Proxy Request in JSON format
func RequestToJSON(req events.APIGatewayProxyRequest) string {
	var buf bytes.Buffer

	body, _ := json.Marshal(req)
	json.HTMLEscape(&buf, body)

	return buf.String()
}
