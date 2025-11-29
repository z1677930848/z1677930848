// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package iplibrary

import "github.com/TeaOSLab/EdgeCommon/pkg/nodeutils"

type Encrypt struct {
}

func NewEncrypt() *Encrypt {
	return &Encrypt{}
}

func (this *Encrypt) Encode(srcData []byte, password string) ([]byte, error) {
	var method = nodeutils.AES256CFBMethod{}
	err := method.Init([]byte(password), []byte(password))
	if err != nil {
		return nil, err
	}

	return method.Encrypt(srcData)
}

func (this *Encrypt) Decode(encodedData []byte, password string) ([]byte, error) {
	var method = nodeutils.AES256CFBMethod{}
	err := method.Init([]byte(password), []byte(password))
	if err != nil {
		return nil, err
	}

	return method.Decrypt(encodedData)
}
