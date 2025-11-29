// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package nodeconfigs

type NodeLevel struct {
	Name        string `yaml:"name" json:"name"`
	Code        int    `yaml:"code" json:"code"`
	Description string `yaml:"description" json:"description"`
}

func FindAllNodeLevels() []*NodeLevel {
	return []*NodeLevel{
		{Name: "Edge", Code: 1, Description: "standard edge node"},
		{Name: "L2", Code: 2, Description: "edge node with upstream backhaul"},
	}
}

func FindNodeLevel(level int) *NodeLevel {
	levels := FindAllNodeLevels()
	idx := level - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(levels) {
		idx = 0
	}
	return levels[idx]
}
