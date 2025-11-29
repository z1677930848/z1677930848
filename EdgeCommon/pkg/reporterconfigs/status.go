// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package reporterconfigs

type Status struct {
	IP               string `json:"ip"`
	OS               string `json:"os"`
	OSName           string `json:"osName"`
	Username         string `json:"username"`
	BuildVersion     string `json:"buildVersion"`     // 缂栬瘧鐗堟湰
	BuildVersionCode uint32 `json:"buildVersionCode"` // 鐗堟湰鏁板瓧
	UpdatedAt        int64  `json:"updatedAt"`        // 鏇存柊鏃堕棿

	Location string `json:"location"` // 浠嶪P鏌ヨ鍒扮殑Location
	ISP      string `json:"isp"`      // 浠嶪P鏌ヨ鍒扮殑ISP
}
