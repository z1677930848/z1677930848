// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

// AccessLogESStorageConfig ElasticSearch瀛樺偍绛栫暐
type AccessLogESStorageConfig struct {
	Endpoint     string `yaml:"endpoint" json:"endpoint"`
	Index        string `yaml:"index" json:"index"`
	MappingType  string `yaml:"mappingType" json:"mappingType"`
	Username     string `yaml:"username" json:"username"`
	Password     string `yaml:"password" json:"password"`
	IsDataStream bool   `yaml:"isDataStream" json:"isDataStream"` // 鏄惁涓篋ata Stream妯″紡
}
