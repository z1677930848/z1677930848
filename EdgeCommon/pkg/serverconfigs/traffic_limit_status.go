// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import timeutil "github.com/iwind/TeaGo/utils/time"

type TrafficLimitTarget = string

const (
	TrafficLimitTargetTraffic              TrafficLimitTarget = "traffic"
	TrafficLimitTargetRequest              TrafficLimitTarget = "request"
	TrafficLimitTargetWebsocketConnections TrafficLimitTarget = "websocketConnections"
)

// TrafficLimitStatus 娴侀噺闄愬埗鐘舵€?
type TrafficLimitStatus struct {
	UntilDay   string `yaml:"untilDay" json:"untilDay"`     // 鏈夋晥鏃ユ湡锛屾牸寮廦YYYMMDD
	PlanId     int64  `yaml:"planId" json:"planId"`         // 濂楅ID
	DateType   string `yaml:"dateType" json:"dateType"`     // 鏃ユ湡绫诲瀷 day|month
	TargetType string `yaml:"targetType" json:"targetType"` // 闄愬埗绫诲瀷锛歵raffic|request|...
}

func (this *TrafficLimitStatus) IsValid() bool {
	if len(this.UntilDay) == 0 {
		return false
	}
	return this.UntilDay >= timeutil.Format("Ymd")
}
