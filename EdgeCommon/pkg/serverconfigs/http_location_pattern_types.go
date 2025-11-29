package serverconfigs

import "github.com/iwind/TeaGo/maps"

type HTTPLocationPatternType = int

const (
	HTTPLocationPatternTypePrefix HTTPLocationPatternType = 1
	HTTPLocationPatternTypeSuffix HTTPLocationPatternType = 4
	HTTPLocationPatternTypeExact  HTTPLocationPatternType = 2
	HTTPLocationPatternTypeRegexp HTTPLocationPatternType = 3
)

func AllLocationPatternTypes() []maps.Map {
	return []maps.Map{
		{"name": "prefix", "type": HTTPLocationPatternTypePrefix, "description": "match path prefix"},
		{"name": "suffix", "type": HTTPLocationPatternTypeSuffix, "description": "match path suffix"},
		{"name": "exact", "type": HTTPLocationPatternTypeExact, "description": "exact path match"},
		{"name": "regexp", "type": HTTPLocationPatternTypeRegexp, "description": "match using regular expression"},
	}
}

func FindLocationPatternType(patternType int) maps.Map {
	for _, t := range AllLocationPatternTypes() {
		if t["type"] == patternType {
			return t
		}
	}
	return nil
}

func FindLocationPatternTypeName(patternType int) string {
	t := FindLocationPatternType(patternType)
	if t == nil {
		return ""
	}
	return t["name"].(string)
}
