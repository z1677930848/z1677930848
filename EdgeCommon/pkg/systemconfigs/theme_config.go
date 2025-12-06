// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package systemconfigs

func DefaultThemeBackgroundColors() []string {
	return []string{
		"14539A",
		"276AC6",
		"0081AF",
		"473BF0",
		"ACADBC",
		"9B9ECE",
		"C96480",
		"B47978",
		"B1AE91",
		"49A078",
		"46237A",
		"000500",
	}
}

// ThemeConfig 椋庢牸妯℃澘璁剧疆
type ThemeConfig struct {
	BackgroundColor string `yaml:"backgroundColor" json:"backgroundColor"` // 鑳屾櫙鑹诧紝16杩涘埗锛屼笉闇€瑕佸甫浜曞彿锛?锛?
}
