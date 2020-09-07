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
			Preferred:    "LMHD",
			Canonical:    "LTTSD",
			Alternatives: []string{},
		},
		FullName: PCA{
			Preferred: "Lucy Top Tier Shitposting Davinhart",
			Canonical: "Lucy Top Tier Shitposting Davinhart",
			Alternatives: []string{
				"Lucy Top Tier Shitposting Davies",
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
				Preferred:    "Top",
				Canonical:    "Top",
				Alternatives: []string{},
			},
			PCA{
				Preferred:    "Tier",
				Canonical:    "Tier",
				Alternatives: []string{},
			},
			PCA{
				Preferred:    "Shitposting",
				Canonical:    "Shitposting",
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
