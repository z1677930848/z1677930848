// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package systemconfigs

type BandwidthUnit = string

const (
	BandwidthUnitByte BandwidthUnit = "byte"
	BandwidthUnitBit  BandwidthUnit = "bit"
)

type BandwidthAlgo = string // 甯﹀绠楁硶

const (
	BandwidthAlgoSecondly BandwidthAlgo = "secondly" // 鎸夌绠?
	BandwidthAlgoAvg      BandwidthAlgo = "avg"      // N鍒嗛挓骞冲潎
)
