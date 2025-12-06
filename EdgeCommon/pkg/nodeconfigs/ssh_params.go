// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package nodeconfigs

func DefaultSSHParams() *SSHParams {
	return &SSHParams{Port: 22}
}

type SSHParams struct {
	Port int `json:"port"`
}
