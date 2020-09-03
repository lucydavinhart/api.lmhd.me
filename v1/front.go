package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rodaine/hclencoder"
)

// PKFront corresponds to a response from the PK fronters API
// https://pluralkit.me/api/#get-s-id-fronters
type PKFront struct {
	Timestamp time.Time `json:"timestamp"`
	Members   []struct {
		ID          string      `json:"id"`
		Name        string      `json:"name"`
		Color       string      `json:"color"`
		DisplayName string      `json:"display_name"`
		Birthday    interface{} `json:"birthday"`
		Pronouns    string      `json:"pronouns"`
		AvatarURL   string      `json:"avatar_url"`
		Description interface{} `json:"description"`
		ProxyTags   []struct {
			Prefix string      `json:"prefix"`
			Suffix interface{} `json:"suffix"`
		} `json:"proxy_tags"`
		KeepProxy          bool        `json:"keep_proxy"`
		Privacy            interface{} `json:"privacy"`
		Visibility         interface{} `json:"visibility"`
		NamePrivacy        interface{} `json:"name_privacy"`
		DescriptionPrivacy interface{} `json:"description_privacy"`
		BirthdayPrivacy    interface{} `json:"birthday_privacy"`
		PronounPrivacy     interface{} `json:"pronoun_privacy"`
		AvatarPrivacy      interface{} `json:"avatar_privacy"`
		MetadataPrivacy    interface{} `json:"metadata_privacy"`
		Created            time.Time   `json:"created"`
		Prefix             string      `json:"prefix"`
		Suffix             interface{} `json:"suffix"`
	} `json:"members"`
}

// GetFront requests the fronter from PluralKit
func GetFront() PKFront {
	url := "https://api.pluralkit.me/v1/s/" +
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

	var front PKFront
	err = json.Unmarshal(bodyBytes, &front)
	if err != nil {
		log.Fatal(err)
	}

	return front
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
