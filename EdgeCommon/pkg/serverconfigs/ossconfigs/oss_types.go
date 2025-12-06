// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package ossconfigs

import "errors"

type OSSType = string

type OSSTypeDefinition struct {
	Name             string `json:"name"`
	Code             string `json:"code"`
	BucketOptionName string `json:"bucketOptionName"`
	BucketIgnored    bool   `json:"bucketIgnored"` // 鏄惁蹇界暐Bucket鍚嶇О
}

func FindAllOSSTypes() []*OSSTypeDefinition {
	return []*OSSTypeDefinition{}
}

func FindOSSType(code string) *OSSTypeDefinition {
	for _, t := range FindAllOSSTypes() {
		if t.Code == code {
			return t
		}
	}
	return nil
}

func DecodeOSSOptions(ossType OSSType, optionsJSON []byte) (any, error) {
	return nil, errors.New("'" + ossType + "' has not been supported")
}
