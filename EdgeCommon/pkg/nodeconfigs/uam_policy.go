// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package nodeconfigs

func NewUAMPolicy() *UAMPolicy {
	return &UAMPolicy{}
}

type UAMPolicy struct {
	IsOn bool `yaml:"isOn" json:"isOn"`
}

func (this *UAMPolicy) Init() error {
	return nil
}
