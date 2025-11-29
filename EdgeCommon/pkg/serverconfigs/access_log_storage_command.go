// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

// AccessLogCommandStorageConfig 閫氳繃鍛戒护琛屽瓨鍌?
type AccessLogCommandStorageConfig struct {
	Command string   `yaml:"command" json:"command"`
	Args    []string `yaml:"args" json:"args"`
	Dir     string   `yaml:"dir" json:"dir"`
}
