// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package reporterconfigs

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
)

type ReportLevel = string

const (
	ReportLevelGood   ReportLevel = "good"
	ReportLevelNormal ReportLevel = "normal"
	ReportLevelBad    ReportLevel = "bad"
	ReportLevelBroken ReportLevel = "broken"
)

func FindAllReportLevels() []*shared.Definition {
	return []*shared.Definition{
		{
			Name: "鑹ソ",
			Code: ReportLevelGood,
		},
		{
			Name: "姝ｅ父",
			Code: ReportLevelNormal,
		},
		{
			Name: "涓嶈壇",
			Code: ReportLevelBad,
		},
		{
			Name: "閿欒",
			Code: ReportLevelBroken,
		},
	}
}

func FindReportLevelName(level ReportLevel) string {
	for _, def := range FindAllReportLevels() {
		if def.Code == level {
			return def.Name
		}
	}
	return ""
}
