// Stub types for DNS.LA provider (community edition)
package dnsla

import "errors"

// BaseResponse base response structure
type BaseResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (this *BaseResponse) Success() bool {
	return this.Status == "success" || this.Code == 200
}

func (this *BaseResponse) Error() error {
	if this.Success() {
		return nil
	}
	return errors.New(this.Message)
}

// DomainListResponse domain list response
type DomainListResponse struct {
	BaseResponse
	Data struct {
		Results []struct {
			Domain string `json:"domain"`
		} `json:"results"`
	} `json:"data"`
}

// RecordListResponse record list response
type RecordListResponse struct {
	BaseResponse
	Data struct {
		Results []struct {
			Id       string `json:"id"`
			Host     string `json:"host"`
			Type     int    `json:"type"`
			Data     string `json:"data"`
			TTL      int    `json:"ttl"`
			LineCode string `json:"lineCode"`
		} `json:"results"`
	} `json:"data"`
}

// AllLineListResponse all line list response
type AllLineListResponse struct {
	BaseResponse
	Data []AllLineListResponseChild `json:"data"`
}

type AllLineListResponseChild struct {
	Id       string                     `json:"id"`
	Name     string                     `json:"name"`
	Code     string                     `json:"code"`
	Children []AllLineListResponseChild `json:"children"`
}

// RecordCreateResponse record create response
type RecordCreateResponse struct {
	BaseResponse
	Data struct {
		Id string `json:"id"`
	} `json:"data"`
}

// RecordUpdateResponse record update response
type RecordUpdateResponse struct {
	BaseResponse
}

// RecordDeleteResponse record delete response
type RecordDeleteResponse struct {
	BaseResponse
}

// DomainResponse domain response
type DomainResponse struct {
	BaseResponse
	Data struct {
		Id string `json:"id"`
	} `json:"data"`
}
