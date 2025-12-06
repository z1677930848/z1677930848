// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import (
	"strings"

	stringutil "github.com/iwind/TeaGo/utils/string"
)

type ScriptConfig struct {
	IsPrior bool `yaml:"isPrior" json:"isPrior"`
	IsOn    bool `yaml:"isOn" json:"isOn"`

	Code            string `yaml:"code" json:"code"`                       // 褰撳墠杩愯鐨勪唬鐮?
	AuditingCode    string `yaml:"auditingCode" json:"auditingCode"`       // 瀹℃牳涓殑浠ｇ爜
	AuditingCodeMD5 string `yaml:"auditingCodeMD5" json:"auditingCodeMD5"` // 瀹℃牳涓殑浠ｇ爜MD5

	realCode string
}

func (this *ScriptConfig) Init() error {
	this.realCode = this.TrimCode()

	return nil
}

func (this *ScriptConfig) TrimCode() string {
	return strings.TrimSpace(this.Code)
}

func (this *ScriptConfig) RealCode() string {
	return this.realCode
}

func (this *ScriptConfig) CodeMD5() string {
	return stringutil.Md5(this.TrimCode())
}
