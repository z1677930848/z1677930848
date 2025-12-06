// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package ossconfigs

import "strings"

func IsOSSProtocol(protocol string) bool {
	return strings.HasPrefix(protocol, "oss:")
}
