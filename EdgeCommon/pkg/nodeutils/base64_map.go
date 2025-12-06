// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package nodeutils

import (
	"encoding/base64"
	"encoding/json"
	"github.com/iwind/TeaGo/maps"
)

// Base64EncodeMap 瀵筂ap杩涜Base64缂栫爜
func Base64EncodeMap(m maps.Map) (string, error) {
	if m == nil {
		m = maps.Map{}
	}
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	var result = base64.StdEncoding.EncodeToString(data)
	return result, nil
}

// Base64DecodeMap 瀵筂ap杩涜Base64瑙ｇ爜
func Base64DecodeMap(encodedString string) (maps.Map, error) {
	data, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return nil, err
	}

	var result = maps.Map{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
