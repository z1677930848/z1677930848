// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package userconfigs

const (
	MaxCacheKeysPerTask int32 = 1000
	MaxCacheKeysPerDay  int32 = 10000
)

type HTTPCacheTaskConfig struct {
	MaxKeysPerTask int32 `yaml:"maxKeysPerTask" json:"maxKeysPerTask"`
	MaxKeysPerDay  int32 `yaml:"maxKeysPerDay" json:"maxKeysPerDay"`
}

func DefaultHTTPCacheTaskConfig() *HTTPCacheTaskConfig {
	return &HTTPCacheTaskConfig{
		MaxKeysPerTask: 0,
		MaxKeysPerDay:  0,
	}
}

// UserServerConfig 鐢ㄦ埛鏈嶅姟璁剧疆
type UserServerConfig struct {
	GroupId                  int64                `yaml:"groupId" json:"groupId"`                                   // 鍒嗙粍
	RequirePlan              bool                 `yaml:"requirePlan" json:"requirePlan"`                           // 蹇呴』浣跨敤濂楅
	EnableStat               bool                 `yaml:"enableStat" json:"enableStat"`                             // 寮€鍚粺璁?
	HTTPCacheTaskPurgeConfig *HTTPCacheTaskConfig `yaml:"httpCacheTaskPurgeConfig" json:"httpCacheTaskPurgeConfig"` // 缂撳瓨浠诲姟鍒犻櫎閰嶇疆
	HTTPCacheTaskFetchConfig *HTTPCacheTaskConfig `yaml:"httpCacheTaskFetchConfig" json:"httpCacheTaskFetchConfig"` // 缂撳瓨浠诲姟棰勭儹閰嶇疆
}

func DefaultUserServerConfig() *UserServerConfig {
	return &UserServerConfig{
		GroupId:                  0,
		RequirePlan:              false,
		EnableStat:               true,
		HTTPCacheTaskPurgeConfig: DefaultHTTPCacheTaskConfig(),
		HTTPCacheTaskFetchConfig: DefaultHTTPCacheTaskConfig(),
	}
}
