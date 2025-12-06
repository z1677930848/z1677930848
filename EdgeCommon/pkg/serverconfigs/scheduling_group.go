// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/schedulingconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/maps"
)

// SchedulingGroup 璐熻浇鍧囪　鍒嗙粍
type SchedulingGroup struct {
	Scheduling *SchedulingConfig `yaml:"scheduling" json:"scheduling"`

	PrimaryOrigins []*OriginConfig
	BackupOrigins  []*OriginConfig

	hasPrimaryOrigins  bool
	hasBackupOrigins   bool
	schedulingIsBackup bool
	schedulingObject   schedulingconfigs.SchedulingInterface
}

// Init 鍒濆鍖?
func (this *SchedulingGroup) Init() error {
	this.hasPrimaryOrigins = len(this.PrimaryOrigins) > 0
	this.hasBackupOrigins = len(this.BackupOrigins) > 0

	if this.Scheduling == nil {
		this.Scheduling = &SchedulingConfig{
			Code:    "random",
			Options: maps.Map{},
		}
	}

	return nil
}

// NextOrigin 鍙栧緱涓嬩竴涓彲鐢ㄦ簮绔?
func (this *SchedulingGroup) NextOrigin(call *shared.RequestCall) *OriginConfig {
	if this.schedulingObject == nil {
		return nil
	}

	if this.Scheduling != nil && call != nil && call.Options != nil {
		for k, v := range this.Scheduling.Options {
			call.Options[k] = v
		}
	}

	var candidate = this.schedulingObject.Next(call)

	// 鏈簡閲嶇疆鐘舵€?
	defer func() {
		if candidate == nil {
			this.schedulingIsBackup = false
		}
	}()

	if candidate == nil {
		// 鍚敤澶囩敤鏈嶅姟鍣?
		if !this.schedulingIsBackup {
			this.SetupScheduling(true, true)
			candidate = this.schedulingObject.Next(call)
			if candidate == nil {
				// 涓嶆鏌ヤ富瑕佹簮绔?
				this.SetupScheduling(false, false)
				candidate = this.schedulingObject.Next(call)
				if candidate == nil {
					// 涓嶆鏌ュ鐢ㄦ簮绔?
					this.SetupScheduling(true, false)
					candidate = this.schedulingObject.Next(call)
					if candidate == nil {
						return nil
					}
				}
			}
		}

		if candidate == nil {
			return nil
		}
	}

	return candidate.(*OriginConfig)
}

// AnyOrigin 鍙栦笅涓€涓换鎰忔簮绔?
func (this *SchedulingGroup) AnyOrigin(excludingOriginIds []int64) *OriginConfig {
	for _, origin := range this.PrimaryOrigins {
		if !origin.IsOn {
			continue
		}
		if !this.containsInt64(excludingOriginIds, origin.Id) {
			return origin
		}
	}
	for _, origin := range this.BackupOrigins {
		if !origin.IsOn {
			continue
		}
		if !this.containsInt64(excludingOriginIds, origin.Id) {
			return origin
		}
	}

	return nil
}

// SetupScheduling 璁剧疆璋冨害绠楁硶
func (this *SchedulingGroup) SetupScheduling(isBackup bool, checkOk bool) {
	// 濡傛灉鍙湁涓€涓簮绔欙紝鍒欏揩閫熻繑鍥烇紝閬垮厤鍥犱负鐘舵€佺殑鏀瑰彉鑰屼笉鍋滃湴杞崲
	if checkOk {
		if len(this.PrimaryOrigins) == 1 && len(this.BackupOrigins) == 0 && this.schedulingObject != nil {
			return
		}
	}

	this.schedulingIsBackup = isBackup

	if this.Scheduling == nil {
		this.schedulingObject = &schedulingconfigs.RandomScheduling{}
	} else {
		typeCode := this.Scheduling.Code
		s := schedulingconfigs.FindSchedulingType(typeCode)
		if s == nil {
			this.Scheduling = nil
			this.schedulingObject = &schedulingconfigs.RandomScheduling{}
		} else {
			this.schedulingObject = s["instance"].(schedulingconfigs.SchedulingInterface)
		}
	}

	if !isBackup {
		for _, origin := range this.PrimaryOrigins {
			if origin.IsOn && (origin.IsOk || !checkOk) {
				this.schedulingObject.Add(origin)
			}
		}
	} else {
		for _, origin := range this.BackupOrigins {
			if origin.IsOn && (origin.IsOk || !checkOk) {
				this.schedulingObject.Add(origin)
			}
		}
	}

	if !this.schedulingObject.HasCandidates() {
		return
	}

	this.schedulingObject.Start()
}

// 鍒ゆ柇鏄惁鍖呭惈int64
func (this *SchedulingGroup) containsInt64(originIds []int64, originId int64) bool {
	for _, id := range originIds {
		if id == originId {
			return true
		}
	}
	return false
}
