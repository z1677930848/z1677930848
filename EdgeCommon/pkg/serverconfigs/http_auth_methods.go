// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

type HTTPAuthType = string

const (
	HTTPAuthTypeBasicAuth  HTTPAuthType = "basicAuth"
	HTTPAuthTypeSubRequest HTTPAuthType = "subRequest"
)

type HTTPAuthTypeDefinition struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

func FindAllHTTPAuthTypes(role string) []*HTTPAuthTypeDefinition {
	return []*HTTPAuthTypeDefinition{
		{Name: "Basic Auth", Code: HTTPAuthTypeBasicAuth, Description: "HTTP Basic authorization"},
		{Name: "Sub Request", Code: HTTPAuthTypeSubRequest, Description: "custom sub-request authentication"},
	}
}
