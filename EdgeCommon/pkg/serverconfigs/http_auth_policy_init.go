// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import (
	"errors"
)

// Init 鍒濆鍖?
func (this *HTTPAuthPolicy) Init() error {
	switch this.Type {
	case HTTPAuthTypeBasicAuth:
		this.method = NewHTTPAuthBasicMethod()
	case HTTPAuthTypeSubRequest:
		this.method = NewHTTPAuthSubRequestMethod()
	}

	if this.method == nil {
		return errors.New("unknown auth method '" + this.Type + "'")
	}
	err := this.method.Init(this.Params)
	if err != nil {
		return err
	}

	return nil
}
