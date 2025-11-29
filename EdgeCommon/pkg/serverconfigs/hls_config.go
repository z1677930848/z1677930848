// Copyright 2024 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package serverconfigs

// HLSConfig HTTP Living Streaming鐩稿叧閰嶇疆
type HLSConfig struct {
	IsPrior    bool                 `yaml:"isPrior" json:"isPrior"`
	Encrypting *HLSEncryptingConfig `yaml:"encrypting" json:"encrypting"` // 鍔犲瘑璁剧疆
}

func (this *HLSConfig) Init() error {
	// encrypting
	if this.Encrypting != nil {
		err := this.Encrypting.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *HLSConfig) IsEmpty() bool {
	if this.Encrypting != nil && this.Encrypting.IsOn {
		return false
	}

	return true
}
