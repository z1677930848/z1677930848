// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

// AccessLogQueueConfig 璁块棶鏃ュ織闃熷垪閰嶇疆
type AccessLogQueueConfig struct {
	MaxLength      int `yaml:"maxLength" json:"maxLength"`           // 闃熷垪鏈€澶ч暱搴?
	CountPerSecond int `yaml:"countPerSecond" json:"countPerSecond"` // 姣忕鍐欏叆鏁伴噺
	Percent        int `yaml:"percent" json:"percent"`               // 姣斾緥锛屽鏋滀负0-100锛岄粯璁や负100

	EnableAutoPartial bool  `yaml:"enableAutoPartial" json:"enableAutoPartial"` // 鏄惁鍚姩鑷姩鍒嗚〃
	RowsPerTable      int64 `yaml:"rowsPerTable" json:"rowsPerTable"`           // 鍗曡〃鏈€澶ц鏁?
}
