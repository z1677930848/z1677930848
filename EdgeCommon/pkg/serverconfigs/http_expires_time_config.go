// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
)

// HTTPExpiresTimeConfig 鍙戦€佸埌瀹㈡埛绔殑杩囨湡鏃堕棿璁剧疆
type HTTPExpiresTimeConfig struct {
	IsPrior       bool                 `yaml:"isPrior" json:"isPrior"`             // 鏄惁瑕嗙洊鐖剁骇璁剧疆
	IsOn          bool                 `yaml:"isOn" json:"isOn"`                   // 鏄惁鍚敤
	Overwrite     bool                 `yaml:"overwrite" json:"overwrite"`         // 鏄惁瑕嗙洊
	AutoCalculate bool                 `yaml:"autoCalculate" json:"autoCalculate"` // 鏄惁鑷姩璁＄畻
	Duration      *shared.TimeDuration `yaml:"duration" json:"duration"`           // 鍛ㄦ湡
}
