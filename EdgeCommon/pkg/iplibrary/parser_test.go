// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package iplibrary_test

import (
	"testing"

	"github.com/TeaOSLab/EdgeCommon/pkg/iplibrary"
)

func TestNewParser(t *testing.T) {
	template, err := iplibrary.NewTemplate("${ipFrom}|${ipTo}|${country}|${any}|${province}|${city}|${provider}")
	if err != nil {
		t.Fatal(err)
	}

	parser, err := iplibrary.NewParser(&iplibrary.ParserConfig{
		Template:    template,
		EmptyValues: []string{"0"},
		Iterator: func(values map[string]string) error {
			t.Log(values)
			return nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	parser.Write([]byte(`0.0.0.0|0.255.255.255|0|0|0|鍐呯綉IP|鍐呯綉IP
1.0.0.0|1.0.0.255|婢冲ぇ鍒╀簹|0|0|0|0
1.0.1.0|1.0.3.255|涓浗|0|绂忓缓鐪亅绂忓窞甯倈鐢典俊
1.0.4.0|1.0.7.255|婢冲ぇ鍒╀簹|0|缁村鍒╀簹|澧ㄥ皵鏈瑋0
1.0.8.0|1.0.15.255|涓浗|0|骞夸笢鐪亅骞垮窞甯倈鐢典俊
1.0.16.0|1.0.31.255|鏃ユ湰|0|0|0|0
1.0.32.0|1.0.63.255|涓浗|0|骞夸笢鐪亅骞垮窞甯倈鐢典俊
1.0.64.0|1.0.79.255|鏃ユ湰|0|骞垮矝鍘縷0|0
1.0.80.0|1.0.127.255|鏃ユ湰|0|鍐堝北鍘縷0|0
1.0.128.0|1.0.128.255|娉板浗|0|娓呰幈搴渱0|TOT
1.0.129.0|1.0.132.191|娉板浗|0|鏇艰胺|鏇艰胺|TOT`))

	err = parser.Parse()
	if err != nil {
		t.Fatal(err)
	}
}
