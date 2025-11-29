// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.
//
// Rewritten to provide a minimal but valid description list for
// IP address related threshold items used inside the node console.

package nodeconfigs

import "github.com/iwind/TeaGo/maps"

// NodeValueOperator 节点值比较操作符
type NodeValueOperator = string

const (
	NodeValueOperatorGt  NodeValueOperator = "gt"  // 大于
	NodeValueOperatorGte NodeValueOperator = "gte" // 大于等于
	NodeValueOperatorLt  NodeValueOperator = "lt"  // 小于
	NodeValueOperatorLte NodeValueOperator = "lte" // 小于等于
	NodeValueOperatorEq  NodeValueOperator = "eq"  // 等于
	NodeValueOperatorNeq NodeValueOperator = "neq" // 不等于
)

// NodeValueDurationUnit 节点值持续时间单位
type NodeValueDurationUnit = string

const (
	NodeValueDurationUnitMinute NodeValueDurationUnit = "minute" // 分钟
	NodeValueDurationUnitHour   NodeValueDurationUnit = "hour"   // 小时
	NodeValueDurationUnitDay    NodeValueDurationUnit = "day"    // 天
)

type IPAddressThresholdItem = string

const (
	IPAddressThresholdItemNodeAvgRequests      IPAddressThresholdItem = "nodeAvgRequests"
	IPAddressThresholdItemNodeAvgTrafficOut    IPAddressThresholdItem = "nodeAvgTrafficOut"
	IPAddressThresholdItemNodeAvgTrafficIn     IPAddressThresholdItem = "nodeAvgTrafficIn"
	IPAddressThresholdItemNodeHealthCheck      IPAddressThresholdItem = "nodeHealthCheck"
	IPAddressThresholdItemGroupAvgRequests     IPAddressThresholdItem = "groupAvgRequests"
	IPAddressThresholdItemGroupAvgTrafficIn    IPAddressThresholdItem = "groupAvgTrafficIn"
	IPAddressThresholdItemGroupAvgTrafficOut   IPAddressThresholdItem = "groupAvgTrafficOut"
	IPAddressThresholdItemClusterAvgRequests   IPAddressThresholdItem = "clusterAvgRequests"
	IPAddressThresholdItemClusterAvgTrafficIn  IPAddressThresholdItem = "clusterAvgTrafficIn"
	IPAddressThresholdItemClusterAvgTrafficOut IPAddressThresholdItem = "clusterAvgTrafficOut"
	IPAddressThresholdItemConnectivity         IPAddressThresholdItem = "connectivity"
)

func FindAllIPAddressThresholdItems() []maps.Map {
	var items = []struct {
		name string
		code IPAddressThresholdItem
		desc string
		unit string
	}{
		{"Node average requests", IPAddressThresholdItemNodeAvgRequests, "Average requests per minute handled by this node.", "req/min"},
		{"Node average downstream traffic", IPAddressThresholdItemNodeAvgTrafficOut, "Average downstream (egress) traffic produced by this node.", "MB"},
		{"Node average upstream traffic", IPAddressThresholdItemNodeAvgTrafficIn, "Average upstream (ingress) traffic received by this node.", "MB"},
		{"Node health check result", IPAddressThresholdItemNodeHealthCheck, "Latest health check status for this node.", ""},
		{"Group average requests", IPAddressThresholdItemGroupAvgRequests, "Average requests per minute for the node group.", "req/min"},
		{"Group average downstream traffic", IPAddressThresholdItemGroupAvgTrafficOut, "Average downstream traffic for the node group.", "MB"},
		{"Group average upstream traffic", IPAddressThresholdItemGroupAvgTrafficIn, "Average upstream traffic for the node group.", "MB"},
		{"Cluster average requests", IPAddressThresholdItemClusterAvgRequests, "Average requests per minute for the cluster.", "req/min"},
		{"Cluster average downstream traffic", IPAddressThresholdItemClusterAvgTrafficOut, "Average downstream traffic for the cluster.", "MB"},
		{"Cluster average upstream traffic", IPAddressThresholdItemClusterAvgTrafficIn, "Average upstream traffic for the cluster.", "MB"},
		{"Connectivity", IPAddressThresholdItemConnectivity, "Connectivity score (0-100) collected from regional monitors.", "%"},
	}

	result := make([]maps.Map, 0, len(items))
	for _, item := range items {
		result = append(result, maps.Map{
			"name":        item.name,
			"code":        item.code,
			"description": item.desc,
			"unit":        item.unit,
		})
	}
	return result
}

type IPAddressThresholdConfig struct {
	Id      int64                             `json:"id"`
	Items   []*IPAddressThresholdItemConfig   `json:"items"`
	Actions []*IPAddressThresholdActionConfig `json:"actions"`
}

type IPAddressThresholdItemConfig struct {
	Item         IPAddressThresholdItem `json:"item"`
	Operator     NodeValueOperator      `json:"operator"`
	Value        float64                `json:"value"`
	Duration     int                    `json:"duration"`
	DurationUnit NodeValueDurationUnit  `json:"durationUnit"`
	Options      maps.Map               `json:"options"`
}

type IPAddressThresholdActionConfig struct {
	Action  string   `json:"action"`
	Options maps.Map `json:"options"`
}

type IPAddressThresholdAction = string

const (
	IPAddressThresholdActionUp      IPAddressThresholdAction = "up"
	IPAddressThresholdActionDown    IPAddressThresholdAction = "down"
	IPAddressThresholdActionNotify  IPAddressThresholdAction = "notify"
	IPAddressThresholdActionSwitch  IPAddressThresholdAction = "switch"
	IPAddressThresholdActionWebHook IPAddressThresholdAction = "webHook"
)

func FindAllIPAddressThresholdActions() []maps.Map {
	return []maps.Map{
		{"name": "Up", "code": IPAddressThresholdActionUp, "description": "IP up"},
		{"name": "Down", "code": IPAddressThresholdActionDown, "description": "IP down"},
		{"name": "Notify", "code": IPAddressThresholdActionNotify, "description": "Notify"},
		{"name": "Switch", "code": IPAddressThresholdActionSwitch, "description": "Switch IP"},
		{"name": "WebHook", "code": IPAddressThresholdActionWebHook, "description": "WebHook"},
	}
}
