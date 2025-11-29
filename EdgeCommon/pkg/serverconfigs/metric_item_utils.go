// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import "github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"

type MetricItemCategory = string

type MetricItemPeriodUnit = string

type MetricItemValueType = string

const (
	MetricItemCategoryHTTP MetricItemCategory = "http"
	MetricItemCategoryTCP  MetricItemCategory = "tcp"
	MetricItemCategoryUDP  MetricItemCategory = "udp"
)

const (
	MetricItemPeriodUnitMinute MetricItemPeriodUnit = "minute"
	MetricItemPeriodUnitHour   MetricItemPeriodUnit = "hour"
	MetricItemPeriodUnitDay    MetricItemPeriodUnit = "day"
	MetricItemPeriodUnitWeek   MetricItemPeriodUnit = "week"
	MetricItemPeriodUnitMonth  MetricItemPeriodUnit = "month"
)

const (
	MetricItemValueTypeCount MetricItemValueType = "count"
	MetricItemValueTypeByte  MetricItemValueType = "byte"
)

func FindAllMetricItemCategoryCodes() []MetricItemCategory {
	return []MetricItemCategory{MetricItemCategoryHTTP, MetricItemCategoryTCP, MetricItemCategoryUDP}
}

func FindAllMetricKeyDefinitions(category MetricItemCategory) []*shared.Definition {
	switch category {
	case MetricItemCategoryHTTP:
		return []*shared.Definition{
			{Name: "client address", Code: "${remoteAddr}", Description: "client IP from headers"},
			{Name: "raw client address", Code: "${rawRemoteAddr}", Description: "direct connection IP"},
			{Name: "username", Code: "${remoteUser}", Description: "basic auth username"},
			{Name: "request URI", Code: "${requestURI}", Description: "path with query"},
			{Name: "request path", Code: "${requestPath}", Description: "path only"},
			{Name: "request URL", Code: "${requestURL}", Description: "full url"},
			{Name: "method", Code: "${requestMethod}", Description: "HTTP method"},
			{Name: "scheme", Code: "${scheme}", Description: "http or https"},
			{Name: "extension", Code: "${requestPathExtension}", Description: "file extension"},
			{Name: "lower extension", Code: "${requestPathLowerExtension}", Description: "lower case extension"},
			{Name: "host", Code: "${host}", Description: "host header"},
			{Name: "protocol", Code: "${proto}", Description: "HTTP protocol"},
			{Name: "arg", Code: "${arg.NAME}", Description: "single query arg"},
			{Name: "referer", Code: "${referer}", Description: "referer url"},
			{Name: "referer host", Code: "${referer.host}", Description: "referer host"},
			{Name: "header", Code: "${header.NAME}", Description: "single header"},
			{Name: "cookie", Code: "${cookie.NAME}", Description: "single cookie"},
			{Name: "status", Code: "${status}", Description: "response status"},
			{Name: "response content type", Code: "${response.contentType}", Description: "response content type"},
		}
	}
	return []*shared.Definition{}
}

type MetricValueDefinition struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Type        string `json:"type"`
}

func FindAllMetricValueDefinitions(category MetricItemCategory) []*MetricValueDefinition {
	switch category {
	case MetricItemCategoryHTTP, MetricItemCategoryTCP, MetricItemCategoryUDP:
		return []*MetricValueDefinition{
			{Name: "requests", Code: "${countRequest}", Type: MetricItemValueTypeCount},
			{Name: "connections", Code: "${countConnection}", Type: MetricItemValueTypeCount},
			{Name: "traffic out", Code: "${countTrafficOut}", Type: MetricItemValueTypeByte},
			{Name: "traffic in", Code: "${countTrafficIn}", Type: MetricItemValueTypeByte},
		}
	}
	return []*MetricValueDefinition{}
}

func FindMetricValueName(category MetricItemCategory, code string) string {
	for _, def := range FindAllMetricValueDefinitions(category) {
		if def.Code == code {
			return def.Name
		}
	}
	return code
}

func FindMetricValueType(category MetricItemCategory, code string) string {
	for _, def := range FindAllMetricValueDefinitions(category) {
		if def.Code == code {
			return def.Type
		}
	}
	return MetricItemValueTypeCount
}

func HumanMetricTime(periodUnit MetricItemPeriodUnit, value string) string {
	switch periodUnit {
	case MetricItemPeriodUnitMonth:
		if len(value) == 6 {
			return value[:4] + "-" + value[4:]
		}
	case MetricItemPeriodUnitWeek:
		if len(value) == 6 {
			return value[:4] + "-" + value[4:]
		}
	case MetricItemPeriodUnitDay:
		if len(value) == 8 {
			return value[:4] + "-" + value[4:6] + "-" + value[6:]
		}
	case MetricItemPeriodUnitHour:
		if len(value) == 10 {
			return value[:4] + "-" + value[4:6] + "-" + value[6:8] + " " + value[8:]
		}
	case MetricItemPeriodUnitMinute:
		if len(value) == 12 {
			return value[:4] + "-" + value[4:6] + "-" + value[6:8] + " " + value[8:10] + ":" + value[10:]
		}
	}
	return value
}

func FindMetricPeriodUnitName(unit string) string {
	switch unit {
	case MetricItemPeriodUnitMonth:
		return "month"
	case MetricItemPeriodUnitWeek:
		return "week"
	case MetricItemPeriodUnitDay:
		return "day"
	case MetricItemPeriodUnitHour:
		return "hour"
	case MetricItemPeriodUnitMinute:
		return "minute"
	}
	return ""
}
