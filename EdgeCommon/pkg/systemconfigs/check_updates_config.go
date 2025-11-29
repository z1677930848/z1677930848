// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package systemconfigs

// CheckUpdatesConfig 妫€鏌ユ洿鏂伴厤缃?
type CheckUpdatesConfig struct {
	AutoCheck      bool   `yaml:"autoCheck" json:"autoCheck"`           // 鏄惁寮€鍚嚜鍔ㄦ鏌?
	IgnoredVersion string `yaml:"ignoredVersion" json:"ignoredVersion"` // 涓婃蹇界暐鐨勭増鏈?
}

func NewCheckUpdatesConfig() *CheckUpdatesConfig {
	return &CheckUpdatesConfig{
		AutoCheck: true,
	}
}
