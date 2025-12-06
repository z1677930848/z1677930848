// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package nodeconfigs

import "github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"

// HTTPPagesPolicy 鍏ㄥ眬鐨凥TTP鑷畾涔夐〉闈㈣缃?
type HTTPPagesPolicy struct {
	IsOn  bool                            `json:"isOn" yaml:"isOn"`   // 鏄惁鍚敤
	Pages []*serverconfigs.HTTPPageConfig `json:"pages" yaml:"pages"` // 鑷畾涔夐〉闈?
}

func NewHTTPPagesPolicy() *HTTPPagesPolicy {
	return &HTTPPagesPolicy{}
}

func (this *HTTPPagesPolicy) Init() error {
	if len(this.Pages) > 0 {
		for _, page := range this.Pages {
			err := page.Init()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
