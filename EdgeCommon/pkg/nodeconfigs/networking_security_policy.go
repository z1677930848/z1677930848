// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package nodeconfigs

import (
	"bytes"
	"encoding/json"
)

type NetworkSecurityStatus = string

const (
	NetworkSecurityStatusAuto NetworkSecurityStatus = "auto"
	NetworkSecurityStatusOn   NetworkSecurityStatus = "on"
	NetworkSecurityStatusOff  NetworkSecurityStatus = "off"
)

// NetworkSecurityPolicy 鑺傜偣缃戠粶瀹夊叏绛栫暐
type NetworkSecurityPolicy struct {
	Status NetworkSecurityStatus `json:"status"` // 鍚敤鐘舵€?

	TCP  struct{} `json:"tcp"`  // TODO
	UDP  struct{} `json:"udp"`  // TODO
	ICMP struct{} `json:"icmp"` // TODO
}

func NewNetworkSecurityPolicy() *NetworkSecurityPolicy {
	var policy = &NetworkSecurityPolicy{}
	policy.Status = NetworkSecurityStatusAuto
	return policy
}

func (this *NetworkSecurityPolicy) Init() error {
	return nil
}

func (this *NetworkSecurityPolicy) AsJSON() ([]byte, error) {
	return json.Marshal(this)
}

func (this *NetworkSecurityPolicy) IsOn() bool {
	return this.Status != NetworkSecurityStatusOff
}

func (this *NetworkSecurityPolicy) IsSame(anotherPolicy *NetworkSecurityPolicy) bool {
	data1, _ := json.Marshal(this)
	data2, _ := json.Marshal(anotherPolicy)
	return bytes.Equal(data1, data2)
}
