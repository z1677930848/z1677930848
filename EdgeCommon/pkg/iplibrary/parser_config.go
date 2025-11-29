// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package iplibrary

type ParserConfig struct {
	Template    *Template
	EmptyValues []string
	Iterator    func(values map[string]string) error
}
