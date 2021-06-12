package main

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rodaine/hclencoder"
	"gopkg.in/yaml.v2"
)

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type Auth struct {
	Authed bool   `json:"authed"`
	Error  *Error `json:"error,omitempty"`
	Scopes Scopes `json:"scopes,omitempty"`
}

func (a Auth) ToJSON() string {
	var buf bytes.Buffer

	body, _ := json.Marshal(a)
	json.HTMLEscape(&buf, body)

	return buf.String()
}

func (a Auth) ToHCL() string {
	hcl, _ := hclencoder.Encode(a)
	return string(hcl)
}

func (a Auth) ToYAML() string {
	yaml, _ := yaml.Marshal(a)
	return string(yaml)
}

// AuthHandler handles requests for /auth etc.
func AuthHandler(req events.APIGatewayProxyRequest, format string) (Response, error) {
	var outputString, outputType, handlerName string
	statusCode := 200

	var auth Auth

	// No required scopes for this API call; just a valid token
	scopes, err := Authorize(req, Scopes{})

	if err != nil {
		auth = Auth{
			Authed: false,
			Error: &Error{
				Message: err.Error(),
				Type:    "AuthError",
			},
		}
		statusCode = 401
	} else {
		// This particular API endpoint will just show the user their API Scopes
		// Other API endpoints will need an Authorization(req, scope) function, which
		// does Authentication and then checks certificate scopes
		auth = Auth{
			Authed: true,
			Scopes: scopes,
		}
	}

	// TODO: got to the point where this is in every function now, huh.
	// We should make an APIResponse interface, return that instead, and have
	// the Handler function in main.go deal with all this

	switch format {
	case "json":
		handlerName = "auth.ToJSON()"
		outputString = auth.ToJSON()
		outputType = "application/json"
	case "yaml":
		handlerName = "auth.ToYAML()"
		outputString = auth.ToYAML()
		outputType = "text/plain; charset=UTF-8"
	default:
		handlerName = "auth.ToHCL()"
		outputString = auth.ToHCL()
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

// Authorize checks if an Authenticated client has the required API Scopes
func Authorize(req events.APIGatewayProxyRequest, requiredScopes Scopes) (Scopes, error) {
	clientScopes, err := Authenticate(req)

	if err != nil {
		return nil, err
	}

	if !requiredScopes.IsAuthorized(clientScopes) {
		return nil, fmt.Errorf("incorrect API scopes. Have: %v, Need: %v", clientScopes, requiredScopes)
	}

	return clientScopes, nil
}

// Authenticate checks for presence of specific request headers and performs Authentication
// If authenticated, returns scopes
func Authenticate(req events.APIGatewayProxyRequest) (scopes Scopes, err error) {
	host, found := req.Headers["Host"]
	if !found {
		return scopes, fmt.Errorf("missing Host header")
	}

	cert, found := req.Headers["x-lmhd-client-cert"]
	if !found {
		return scopes, fmt.Errorf("missing x-lmhd-client-cert header")
	}

	return ValidateClientCertificate(cert, host)
}

// ValidateCertificate checks if a client certificate is valid for the current
// instance of the API.
// If it is, then return the certificate scopes.
//
// This uses an x.509 PUBLIC certificate. The API does not make use of the
// private certificate (yet), due to the API clients being iOS Shortcuts
//
// This is something which would be useful to change in future
func ValidateClientCertificate(certRawBase64, host string) (Scopes, error) {
	// Which CAs do we trust?
	var trustedRootCAs = []string{
		// TODO: Parameterise this with an env var
		"https://vault.lmhd.me/v1/pki/root/ca",
	}

	// Get root cert from our trusted CA list
	roots := x509.NewCertPool()
	for _, caURL := range trustedRootCAs {
		root, err := GetCAFromURL(caURL)
		if err != nil {
			return []string{}, fmt.Errorf("failed to parse inter certificate: %v", err)
		}
		roots.AddCert(root)
	}

	// Parse our certificate string

	certRaw, err := base64.URLEncoding.DecodeString(certRawBase64)
	if err != nil {
		return []string{}, fmt.Errorf("could not base64 decode certificate")
	}

	block, _ := pem.Decode(certRaw)
	if block == nil {
		return []string{}, fmt.Errorf("could not decode PEM: nil block")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return []string{}, fmt.Errorf("could not parse certificate: %v", err)
	}

	// Get Intermediary cert from main cert
	//
	// This works on the assumption of a chain like:
	//   Root --> Inter --> Cert
	//
	// If there are any further intermediaries in the chain, this will not work
	//
	// We could also require the client present its Issuing CA as an HTTP header
	// but that is out-of-scope of this proof-of-concept
	inters := x509.NewCertPool()
	for _, caURL := range cert.IssuingCertificateURL {
		inter, err := GetCAFromURL(caURL)
		if err != nil {
			return []string{}, fmt.Errorf("failed to parse inter certificate: %v", err)
		}
		inters.AddCert(inter)
	}

	// Verify

	opts := x509.VerifyOptions{
		// Our trusted CA
		Roots: roots,

		// The issuing CA, as reported by the cert
		Intermediates: inters,

		// Validate this is a client certificate
		KeyUsages: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
		},
	}

	if _, err := cert.Verify(opts); err != nil {
		return []string{}, fmt.Errorf("failed to verify certificate: %v", err)
	}

	// Check if our application host matches the OU in the cert
	hostMatch := false
	for _, ou := range cert.Subject.OrganizationalUnit {
		if host == ou {
			hostMatch = true
		}
	}

	if hostMatch {
		return cert.Subject.Organization, nil
	} else {
		return []string{}, fmt.Errorf("cert OUs do not match host")
	}
}

// GetCAFromURL reads a raw Certificate Authority Certificate from a URL and
// returns it as an *x509.Certificate
func GetCAFromURL(address string) (*x509.Certificate, error) {

	resp, err := http.Get(address)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(body)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
