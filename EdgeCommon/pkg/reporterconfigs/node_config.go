// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package reporterconfigs

type NodeConfig struct {
	Id int64 `json:"id"`
}

func (this *NodeConfig) Init() error {
	return nil
}
