package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rodaine/hclencoder"
	"gopkg.in/yaml.v2"
)

type PCA struct {
	Preferred    string   `json:"preferred"`
	Canonical    string   `json:"canonical"`
	Alternatives []string `json:"alternatives"`
}

type Name struct {
	Initials    PCA    `json:"initials"`
	FullName    PCA    `json:"full_name"`
	FirstName   PCA    `json:"first_name"`
	MiddleNames []PCA  `json:"middle_names"`
	FamilyName  PCA    `json:"family_name"`
	Version     string `json:"version"`
}

func GetName() Name {
	return Name{
		Initials: PCA{
			Preferred:    "LMHD",
			Canonical:    "LMHAD",
			Alternatives: []string{},
		},
		FullName: PCA{
			Preferred: "Lucy Mægan Heather Artemis Davinhart",
			Canonical: "Lucy Mægan Heather Artemis Davinhart",
			Alternatives: []string{
				"Lucy Maegan Heather Artemis Davinhart",
				"Lucy Mægan Heather Artemis Davies",
				"Lucy Maegan Heather Artemis Davies",
			},
		},
		FirstName: PCA{
			Preferred: "Lucy",
			Canonical: "Lucy",
			Alternatives: []string{
				"Lucidity",
				"Lusitania",
				"Liùsaidh",
				"Liusaidh",
				"Lucidora",
				"Lucille",
			},
		},
		MiddleNames: []PCA{
			PCA{
				Preferred:    "Mægan",
				Canonical:    "Mægan",
				Alternatives: []string{"Maegan"},
			},
			PCA{
				Preferred:    "Heather",
				Canonical:    "Heather",
				Alternatives: []string{},
			},
			PCA{
				Preferred:    "Artemis",
				Canonical:    "Artemis",
				Alternatives: []string{},
			},
		},
		FamilyName: PCA{
			Preferred:    "Davinhart",
			Canonical:    "Davinhart",
			Alternatives: []string{"Davies"},
		},
		Version: "1.2.0",
	}
}

func (n Name) ToJSON() string {
	var buf bytes.Buffer

	body, _ := json.Marshal(n)
	json.HTMLEscape(&buf, body)

	return buf.String()
}

func (n Name) ToHCL() string {
	hcl, _ := hclencoder.Encode(n)
	return string(hcl)
}

func (n Name) ToYAML() string {
	yaml, _ := yaml.Marshal(n)
	return string(yaml)
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
