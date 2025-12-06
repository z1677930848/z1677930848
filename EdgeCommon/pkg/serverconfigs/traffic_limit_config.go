// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import "github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"

// DefaultTrafficLimitNoticePageBody 杈惧埌娴侀噺闄愬埗鏃堕粯璁ゆ彁绀哄唴瀹?
const DefaultTrafficLimitNoticePageBody = `<!DOCTYPE html>
<html>
<head>
<title>Traffic Limit Exceeded Warning</title>
<body>

<h1>Traffic Limit Exceeded Warning</h1>
<p>The site traffic has exceeded the limit. Please contact with the site administrator.</p>
<address>Request ID: ${requestId}.</address>

</body>
</html>`

// TrafficLimitConfig 娴侀噺闄愬埗
type TrafficLimitConfig struct {
	IsOn bool `yaml:"isOn" json:"isOn"` // 鏄惁鍚敤

	DailySize   *shared.SizeCapacity `yaml:"dailySize" json:"dailySize"`     // 姣忔棩闄愬埗
	MonthlySize *shared.SizeCapacity `yaml:"monthlySize" json:"monthlySize"` // 姣忔湀闄愬埗
	TotalSize   *shared.SizeCapacity `yaml:"totalSize" json:"totalSize"`     // 鎬婚檺鍒?TODO 闇€瑕佸疄鐜?

	NoticePageBody string `yaml:"noticePageBody" json:"noticePageBody"` // 瓒呭嚭闄愬埗鏃剁殑鎻愰啋锛屾敮鎸佽姹傚彉閲?
}

func (this *TrafficLimitConfig) Init() error {
	return nil
}

// DailyBytes 姣忓ぉ闄愬埗
// 涓嶄娇鐢↖nit()鏉ュ垵濮嬪寲鏁版嵁锛屾槸涓轰簡璁╁叾浠栧湴鏂逛笉缁忚繃Init()涔熻兘寰楀埌璁＄畻鍊?
func (this *TrafficLimitConfig) DailyBytes() int64 {
	if this.DailySize != nil {
		return this.DailySize.Bytes()
	}
	return -1
}

// MonthlyBytes 姣忔湀闄愬埗
func (this *TrafficLimitConfig) MonthlyBytes() int64 {
	if this.MonthlySize != nil {
		return this.MonthlySize.Bytes()
	}
	return -1
}

// TotalBytes 鎬婚檺鍒?
func (this *TrafficLimitConfig) TotalBytes() int64 {
	if this.TotalSize != nil {
		return this.TotalSize.Bytes()
	}
	return -1
}

// IsEmpty 妫€鏌ユ槸鍚︽湁闄愬埗鍊?
func (this *TrafficLimitConfig) IsEmpty() bool {
	return !this.IsOn || (this.DailyBytes() <= 0 && this.MonthlyBytes() <= 0 && this.TotalBytes() <= 0)
}

func (this *TrafficLimitConfig) Equals(another *TrafficLimitConfig) bool {
	if another == nil {
		return false
	}

	if this.IsOn != another.IsOn {
		return false
	}

	if !this.equalCapacity(this.DailySize, another.DailySize) {
		return false
	}

	if !this.equalCapacity(this.MonthlySize, another.MonthlySize) {
		return false
	}

	if !this.equalCapacity(this.TotalSize, another.TotalSize) {
		return false
	}

	if this.NoticePageBody != another.NoticePageBody {
		return false
	}

	return true
}

func (this *TrafficLimitConfig) equalCapacity(size1 *shared.SizeCapacity, size2 *shared.SizeCapacity) bool {
	if size1 == size2 { // all are nil
		return true
	}

	if size1 != nil {
		if size2 == nil {
			return false
		}
		return size1.Bytes() == size2.Bytes()
	}

	return false
}
