// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

// HTTPAuthPolicyRef 璁よ瘉绛栫暐寮曠敤
type HTTPAuthPolicyRef struct {
	IsOn         bool            `yaml:"isOn" json:"isOn"`
	AuthPolicyId int64           `yaml:"authPolicyId" json:"authPolicyId"`
	AuthPolicy   *HTTPAuthPolicy `yaml:"authPolicy" json:"authPolicy"`
}

func (this *HTTPAuthPolicyRef) Init() error {
	if this.AuthPolicy != nil {
		err := this.AuthPolicy.Init()
		if err != nil {
			return err
		}
	}
	return nil
}
