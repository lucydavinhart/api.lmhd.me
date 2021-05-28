package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rodaine/hclencoder"
	"gopkg.in/yaml.v2"
)

type Quest struct {
	Quest string `json:"quest"`
}

func (q Quest) ToJSON() string {
	var buf bytes.Buffer

	body, _ := json.Marshal(q)
	json.HTMLEscape(&buf, body)

	return buf.String()
}

func (q Quest) ToHCL() string {
	hcl, _ := hclencoder.Encode(q)
	return string(hcl)
}

func (q Quest) ToYAML() string {
	yaml, _ := yaml.Marshal(q)
	return string(yaml)
}

// QuestHandler handles requests for /quest etc.
func QuestHandler(req events.APIGatewayProxyRequest, format string) (Response, error) {
	var outputString, outputType, handlerName string

	quest := Quest{
		Quest: "Instigator of DevOps Shenanigans, Girl Adjacent Who Go, and Recreational Procrastinatrix",
	}

	switch format {
	case "json":
		handlerName = "quest.ToJSON()"
		outputString = quest.ToJSON()
		outputType = "application/json"
	case "yaml":
		handlerName = "quest.ToYAML()"
		outputString = quest.ToYAML()
		outputType = "text/plain; charset=UTF-8"
	default:
		handlerName = "quest.ToHCL()"
		outputString = quest.ToHCL()
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
