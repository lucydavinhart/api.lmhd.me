package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rodaine/hclencoder"
	"gopkg.in/yaml.v2"
)

// PKFront corresponds to a response from the PK fronters API
// https://pluralkit.me/api/#get-s-id-fronters
type PKFront struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Members   []struct {
		ID          string      `json:"id"`
		UUID        string      `json:"uuid"`
		Name        string      `json:"name"`
		DisplayName interface{} `json:"display_name"`
		Color       string      `json:"color"`
		Birthday    interface{} `json:"birthday"`
		Pronouns    string      `json:"pronouns"`
		AvatarURL   string      `json:"avatar_url"`
		Banner      interface{} `json:"banner"`
		Description interface{} `json:"description"`
		Created     time.Time   `json:"created"`
		KeepProxy   bool        `json:"keep_proxy"`
		ProxyTags   []struct {
			Prefix string      `json:"prefix"`
			Suffix interface{} `json:"suffix"`
		} `json:"proxy_tags"`
		Privacy struct {
			Visibility         string `json:"visibility"`
			NamePrivacy        string `json:"name_privacy"`
			DescriptionPrivacy string `json:"description_privacy"`
			BirthdayPrivacy    string `json:"birthday_privacy"`
			PronounPrivacy     string `json:"pronoun_privacy"`
			AvatarPrivacy      string `json:"avatar_privacy"`
			MetadataPrivacy    string `json:"metadata_privacy"`
		} `json:"privacy"`
	} `json:"members"`
}

// GetFront requests the fronter from PluralKit
func GetFront() PKFront {
	url := "https://api.pluralkit.me/v2/systems/" +
		os.ExpandEnv("${PLURALKIT_SYSTEM_ID}") +
		"/fronters"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", os.ExpandEnv("${PLURALKIT_API_TOKEN}"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response: %s", bodyBytes)

	var front PKFront
	err = json.Unmarshal(bodyBytes, &front)
	if err != nil {
		log.Fatal(err)
	}

	// Assume only one fronter for now
	if front.Members[0].Privacy.Visibility == "public" {
		return front
	} else {
		return PKFront{}
	}
}

// ToJSON outputs PKFront in JSON format
func (f PKFront) ToJSON() string {
	var buf bytes.Buffer

	body, _ := json.Marshal(f)
	json.HTMLEscape(&buf, body)

	return buf.String()
}

// ToHCL outputs PKFront in HCL format
func (f PKFront) ToHCL() string {
	hcl, _ := hclencoder.Encode(f)
	return string(hcl)
}

func (f PKFront) ToYAML() string {
	yaml, _ := yaml.Marshal(f)
	return string(yaml)
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
