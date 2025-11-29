// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import (
	"time"

	"github.com/iwind/TeaGo/Tea"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

// MetricItemConfig 鎸囨爣閰嶇疆
type MetricItemConfig struct {
	Id            int64                `yaml:"id" json:"id"`
	IsOn          bool                 `yaml:"isOn" json:"isOn"`
	Category      MetricItemCategory   `yaml:"category" json:"category"`
	Period        int                  `yaml:"period" json:"period"`         // 缁熻鍛ㄦ湡
	PeriodUnit    MetricItemPeriodUnit `yaml:"periodUnit" json:"periodUnit"` // 缁熻鍛ㄦ湡鍗曚綅
	Keys          []string             `yaml:"keys" json:"keys"`
	Value         string               `yaml:"value" json:"value"`
	Version       int32                `yaml:"version" json:"version"`
	ExpiresPeriod int                  `yaml:"expiresPeriod" json:"expiresPeriod"` // 杩囨湡鍛ㄦ湡

	//sumType                string    // 缁熻绫诲瀷
	baseTime               time.Time // 鍩哄噯鏃堕棿
	hasHTTPConnectionValue bool      // 鏄惁鏈夌粺璁TTP杩炴帴鏁扮殑鏁板€?
}

// Init 鍒濆鍖?
func (this *MetricItemConfig) Init() error {
	// 鎵€鏈夋椂闂翠互 2020-01-01鏃?涓哄熀鍑?
	this.baseTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)

	if this.Period <= 0 {
		this.Period = 1
	}

	if len(this.PeriodUnit) == 0 {
		this.PeriodUnit = MetricItemPeriodUnitDay
	}

	this.hasHTTPConnectionValue = this.Category == MetricItemCategoryHTTP && this.Value == "${countConnection}"

	return nil
}

// CurrentTime 鏍规嵁鍛ㄦ湡璁＄畻鏃堕棿
func (this *MetricItemConfig) CurrentTime() string {
	var t string

	switch this.PeriodUnit {
	case MetricItemPeriodUnitMonth:
		if this.Period > 1 {
			var now = time.Now()
			var months = (now.Year()-this.baseTime.Year())*12 + int(now.Month())
			var delta = months % this.Period
			if delta == 0 {
				t = timeutil.Format("Ym")
			} else {
				t = timeutil.Format("Ym", now.AddDate(0, -delta, 0))
			}
		} else {
			t = timeutil.Format("Ym")
		}
	case MetricItemPeriodUnitWeek:
		if this.Period > 1 {
			var now = time.Now()
			var weeks = int((now.Unix() - this.baseTime.Unix()) / (86400 * 7))
			var delta = weeks % this.Period
			if delta == 0 {
				t = timeutil.Format("YW")
			} else {
				t = timeutil.FormatTime("YW", now.Unix()-int64(delta*7*86400))
			}
		} else {
			t = timeutil.Format("YW")
		}
	case MetricItemPeriodUnitDay:
		if this.Period > 1 {
			var now = time.Now()
			var days = int((now.Unix() - this.baseTime.Unix()) / 86400)
			var delta = days % this.Period
			if delta == 0 {
				t = timeutil.Format("Ymd")
			} else {
				t = timeutil.FormatTime("Ymd", now.Unix()-int64(delta*86400))
			}
		} else {
			t = timeutil.Format("Ymd")
		}
	case MetricItemPeriodUnitHour:
		if this.Period > 1 {
			var now = time.Now()
			var hours = int((now.Unix() - this.baseTime.Unix()) / 3600)
			var delta = hours % this.Period
			if delta == 0 {
				t = timeutil.Format("YmdH")
			} else {
				t = timeutil.FormatTime("YmdH", now.Unix()-int64(delta*3600))
			}
		} else {
			t = timeutil.Format("YmdH")
		}
	case MetricItemPeriodUnitMinute:
		if this.Period > 1 {
			var now = time.Now()
			var minutes = int((now.Unix() - this.baseTime.Unix()) / 60)
			var delta = minutes % this.Period
			if delta == 0 {
				t = timeutil.Format("YmdHi")
			} else {
				t = timeutil.FormatTime("YmdHi", now.Unix()-int64(delta*60))
			}
		} else {
			t = timeutil.Format("YmdHi")
		}
	default:
		return ""
	}
	return t
}

// ServerExpiresTime 鏍规嵁鍛ㄦ湡璁＄畻鏈嶅姟鍣ㄧ鏁版嵁杩囨湡鏃堕棿
func (this *MetricItemConfig) ServerExpiresTime() string {
	switch this.PeriodUnit {
	case MetricItemPeriodUnitMonth:
		return timeutil.Format("Ym", time.Now().AddDate(0, -(this.Period*4), 0))
	case MetricItemPeriodUnitWeek:
		return timeutil.FormatTime("YW", time.Now().Unix()-86400*7*int64(this.Period*5))
	case MetricItemPeriodUnitDay:
		return timeutil.FormatTime("Ymd", time.Now().Unix()-86400*int64(this.Period*32))
	case MetricItemPeriodUnitHour:
		return timeutil.FormatTime("YmdH", time.Now().Unix()-3600*int64(this.Period*25))
	case MetricItemPeriodUnitMinute:
		return timeutil.FormatTime("YmdHi", time.Now().Unix()-60*int64(this.Period*60))
	default:
		return ""
	}
}

// ServerExpiresDay 鏍规嵁鍛ㄦ湡璁＄畻鏈嶅姟鍣ㄧ鏁版嵁杩囨湡鏃ユ湡
func (this *MetricItemConfig) ServerExpiresDay() string {
	switch this.PeriodUnit {
	case MetricItemPeriodUnitMonth:
		if this.ExpiresPeriod > 0 {
			return timeutil.Format("Ymd", time.Now().AddDate(0, -(this.Period*this.ExpiresPeriod), 0))
		} else {
			return timeutil.Format("Ymd", time.Now().AddDate(0, -(this.Period*4), 0))
		}
	case MetricItemPeriodUnitWeek:
		if this.ExpiresPeriod > 0 {
			return timeutil.FormatTime("Ymd", time.Now().Unix()-86400*7*int64(this.Period*this.ExpiresPeriod))
		} else {
			return timeutil.FormatTime("Ymd", time.Now().Unix()-86400*7*int64(this.Period*5))
		}
	case MetricItemPeriodUnitDay:
		if this.ExpiresPeriod > 0 {
			return timeutil.FormatTime("Ymd", time.Now().Unix()-86400*int64(this.Period*this.ExpiresPeriod))
		} else {
			return timeutil.FormatTime("Ymd", time.Now().Unix()-86400*int64(this.Period*8))
		}
	case MetricItemPeriodUnitHour:
		if this.ExpiresPeriod > 0 {
			return timeutil.FormatTime("Ymd", time.Now().Unix()-3600*int64(this.Period*this.ExpiresPeriod)-86400)
		} else {
			return timeutil.FormatTime("Ymd", time.Now().Unix()-3600*int64(this.Period*25)-86400)
		}
	case MetricItemPeriodUnitMinute:
		if this.ExpiresPeriod > 0 {
			return timeutil.FormatTime("Ymd", time.Now().Unix()-60*int64(this.Period*this.ExpiresPeriod)-86400)
		} else {
			return timeutil.FormatTime("Ymd", time.Now().Unix()-60*int64(this.Period*60)-86400)
		}
	default:
		return ""
	}
}

// LocalExpiresTime 鏍规嵁鍛ㄦ湡璁＄畻鏈湴绔繃鏈熸椂闂?
func (this *MetricItemConfig) LocalExpiresTime() string {
	switch this.PeriodUnit {
	case MetricItemPeriodUnitMonth:
		return timeutil.Format("Ym", time.Now().AddDate(0, -(this.Period+1), 0))
	case MetricItemPeriodUnitWeek:
		return timeutil.FormatTime("YW", time.Now().Unix()-86400*7*int64(this.Period+1))
	case MetricItemPeriodUnitDay:
		return timeutil.FormatTime("Ymd", time.Now().Unix()-86400*int64(this.Period+1))
	case MetricItemPeriodUnitHour:
		return timeutil.FormatTime("YmdH", time.Now().Unix()-3600*int64(this.Period+1))
	case MetricItemPeriodUnitMinute:
		return timeutil.FormatTime("YmdHi", time.Now().Unix()-60*int64(this.Period+60))
	default:
		return ""
	}
}

// UploadDuration 涓婁紶鏁版嵁鐨勫懆鏈?
func (this *MetricItemConfig) UploadDuration() time.Duration {
	if Tea.IsTesting() {
		return 5 * time.Second
	}
	switch this.PeriodUnit {
	case MetricItemPeriodUnitMonth:
		return 10 * time.Minute
	case MetricItemPeriodUnitWeek:
		return 10 * time.Minute
	case MetricItemPeriodUnitDay:
		return 10 * time.Minute
	case MetricItemPeriodUnitHour:
		return 5 * time.Minute
	case MetricItemPeriodUnitMinute:
		return time.Duration(this.Period) * time.Minute
	default:
		return 10 * time.Minute
	}
}

func (this *MetricItemConfig) HasHTTPConnectionValue() bool {
	return this.hasHTTPConnectionValue
}
