package main

import (
	"bytes"
	"encoding/json"

	"github.com/rodaine/hclencoder"
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
			Preferred:    "HONCCC",
			Canonical:    "HONK",
			Alternatives: []string{},
		},
		FullName: PCA{
			Preferred:    "Grian Coldwater",
			Canonical:    "Ian Coldwater",
			Alternatives: []string{},
		},
		FirstName: PCA{
			Preferred:    "Grian",
			Canonical:    "Ian",
			Alternatives: []string{},
		},
		MiddleNames: []PCA{},
		FamilyName: PCA{
			Preferred:    "Coldwater",
			Canonical:    "Coldwater",
			Alternatives: []string{},
		},
		Version: "honcc",
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
