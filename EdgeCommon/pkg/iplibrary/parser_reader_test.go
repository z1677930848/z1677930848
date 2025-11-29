// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package iplibrary_test

import (
	"bytes"
	"github.com/TeaOSLab/EdgeCommon/pkg/iplibrary"
	"testing"
)

func TestNewReaderParser(t *testing.T) {
	template, err := iplibrary.NewTemplate("${ipFrom}|${ipTo}|${country}|${any}|${province}|${city}|${provider}")
	if err != nil {
		t.Fatal(err)
	}

	var buf = &bytes.Buffer{}
	buf.WriteString(`8.45.160.0|8.45.161.255|缇庡浗|0|鍗庣洓椤縷瑗块泤鍥緗Level3
8.45.162.0|8.45.162.255|缇庡浗|0|椹惃璇稿|0|Level3
8.45.163.0|8.45.164.255|缇庡浗|0|淇勫嫆鍐坾0|Level3
8.45.165.0|8.45.165.255|缇庡浗|0|鍗庣洓椤縷0|Level3
8.45.166.0|8.45.167.255|缇庡浗|0|鍗庣洓椤縷瑗块泤鍥緗Level3
8.45.168.0|8.127.255.255|缇庡浗|0|0|0|Level3
8.128.0.0|8.128.3.255|涓浗|0|涓婃捣|涓婃捣甯倈闃块噷宸村反
8.128.4.0|8.128.255.255|涓浗|0|0|0|闃块噷宸村反
8.129.0.0|8.129.255.255|涓浗|0|骞夸笢鐪亅娣卞湷甯倈闃块噷浜?
8.130.0.0|8.130.55.255|涓浗|0|鍖椾含|鍖椾含甯倈闃块噷浜?
8.130.56.0|8.131.255.255|涓浗|0|鍖椾含|鍖椾含甯倈闃块噷宸村反
8.132.0.0|8.133.255.255|涓浗|0|涓婃捣|涓婃捣甯倈闃块噷宸村反
8.134.0.0|8.134.20.63|涓浗|0|0|0|闃块噷浜?
8.134.20.64|8.134.79.255|涓浗|0|骞夸笢鐪亅娣卞湷甯倈闃块噷浜?
8.134.80.0|8.191.255.255|涓浗|0|0|0|闃块噷宸村反
8.192.0.0|8.192.0.255|缇庡浗|0|0|0|Level3
8.192.1.0|8.192.1.255|缇庡浗|0|椹惃璇稿|娉㈠＋椤縷Level3
8.192.2.0|8.207.255.255|缇庡浗|0|0|0|Level3
8.208.0.0|8.208.127.255|鑻卞浗|0|浼︽暒|浼︽暒|闃块噷浜?
8.208.128.0|8.208.255.255|鑻卞浗|0|浼︽暒|浼︽暒|闃块噷宸村反
8.209.0.0|8.209.15.255|淇勭綏鏂瘄0|鑾柉绉憒鑾柉绉憒闃块噷浜?
8.209.16.0|8.209.31.255|淇勭綏鏂瘄0|鑾柉绉憒鑾柉绉憒闃块噷宸村反
8.209.32.0|8.209.32.15|涓浗|0|鍙版咕鐪亅鍙板寳|闃块噷浜?
8.209.32.16|8.209.32.255|涓浗|0|鍙版咕鐪亅鍙板寳|闃块噷宸村反
8.209.33.0|8.209.34.255|涓浗|0|鍙版咕鐪亅鍙板寳|闃块噷浜?
8.209.35.0|8.209.35.255|涓浗|0|鍙版咕鐪亅鍙板寳|闃块噷宸村反`)

	var count int
	parser, err := iplibrary.NewReaderParser(buf, &iplibrary.ParserConfig{
		Template:    template,
		EmptyValues: []string{"0"},
		Iterator: func(values map[string]string) error {
			count++
			t.Log(count, values)
			return nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = parser.Parse()
	if err != nil {
		t.Fatal(err)
	}
}
