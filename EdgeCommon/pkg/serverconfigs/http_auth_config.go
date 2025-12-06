// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

// HTTPAuthConfig 璁よ瘉閰嶇疆
type HTTPAuthConfig struct {
	IsPrior    bool                 `yaml:"isPrior" json:"isPrior"`
	IsOn       bool                 `yaml:"isOn" json:"isOn"`
	PolicyRefs []*HTTPAuthPolicyRef `yaml:"policyRefs" json:"policyRefs"`
}

func (this *HTTPAuthConfig) Init() error {
	for _, ref := range this.PolicyRefs {
		err := ref.Init()
		if err != nil {
			return err
		}
	}
	return nil
}
