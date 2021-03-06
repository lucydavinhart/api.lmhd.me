package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rodaine/hclencoder"
	federate "github.com/strawberryutopia/federate-fronter"
	"gopkg.in/yaml.v2"
)

type Federate struct {
	Updated bool   `json:"updated"`
	Fronter string `json:"fronter"`
	Error   *Error `json:"error,omitempty"`
}

func (f Federate) ToJSON() string {
	var buf bytes.Buffer

	body, _ := json.Marshal(f)
	json.HTMLEscape(&buf, body)

	return buf.String()
}

func (f Federate) ToHCL() string {
	hcl, _ := hclencoder.Encode(f)
	return string(hcl)
}

func (f Federate) ToYAML() string {
	yaml, _ := yaml.Marshal(f)
	return string(yaml)
}

// FrontFederateHandler handles requests for /front/federate etc.
func FrontFederateHandler(req events.APIGatewayProxyRequest, format string) (Response, error) {
	var outputString, outputType, handlerName string
	statusCode := 200

	var fed Federate

	// Need scopes to be able to write to this
	_, err := Authorize(req, Scopes{"fronter.federate:write"})

	// TODO: This is kind of a mess.
	// Rather than having lots of nested IFs, do everythign in another function,
	// and return the first error we get

	if err != nil {
		fed = Federate{
			Updated: false,
			Error: &Error{
				Message: err.Error(),
				Type:    "AuthError",
			},
		}
		statusCode = 401
	} else {
		client, err := federate.NewClient()
		if err != nil {
			fed = Federate{
				Updated: false,
				Error: &Error{
					Message: err.Error(),
					Type:    "ClientError",
				},
			}
			statusCode = 401
		} else {
			fronter, err := client.UpdateFromFront()
			if err != nil {
				fed = Federate{
					Updated: false,
					Error: &Error{
						Message: err.Error(),
						Type:    "FederateError",
					},
				}
				statusCode = 401
			} else {
				fed = Federate{
					Updated: true,
					Fronter: fronter,
				}
			}
		}
	}

	// TODO: got to the point where this is in every function now, huh.
	// We should make an APIResponse interface, return that instead, and have
	// the Handler function in main.go deal with all this

	switch format {
	case "json":
		handlerName = "fed.ToJSON()"
		outputString = fed.ToJSON()
		outputType = "application/json"
	case "yaml":
		handlerName = "fed.ToYAML()"
		outputString = fed.ToYAML()
		outputType = "text/plain; charset=UTF-8"
	default:
		handlerName = "fed.ToHCL()"
		outputString = fed.ToHCL()
		outputType = "text/plain; charset=UTF-8"
	}
	fmt.Printf("%v", outputString)

	resp := Response{
		StatusCode:      statusCode,
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
