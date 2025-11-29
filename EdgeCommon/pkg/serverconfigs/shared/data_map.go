// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package shared

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"sync"
)

var dataMapPrefix = []byte("Lingcdn_DATA_MAP:")

// DataMap 浜岃繘鍒舵暟鎹叡浜玀ap
// 鐢ㄦ潵鍑忓皯鐩稿悓鏁版嵁鍗犵敤鐨勭┖闂村拰鍐呭瓨
type DataMap struct {
	Map    map[string][]byte
	locker sync.Mutex
}

// NewDataMap 鏋勫缓瀵硅薄
func NewDataMap() *DataMap {
	return &DataMap{Map: map[string][]byte{}}
}

// Put 鏀惧叆鏁版嵁
func (this *DataMap) Put(data []byte) (keyData []byte) {
	this.locker.Lock()
	defer this.locker.Unlock()
	var key = string(dataMapPrefix) + fmt.Sprintf("%x", md5.Sum(data))
	this.Map[key] = data
	return []byte(key)
}

// Read 璇诲彇鏁版嵁
func (this *DataMap) Read(key []byte) []byte {
	this.locker.Lock()
	defer this.locker.Unlock()
	if bytes.HasPrefix(key, dataMapPrefix) {
		return this.Map[string(key)]
	}
	return key
}
