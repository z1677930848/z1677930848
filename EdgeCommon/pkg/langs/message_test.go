// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package langs

import "testing"

func TestMessageCode_For(t *testing.T) {
	defaultManager.AddLang("en-us").
		Set("name", "Lily")

	var messageCode MessageCode = "name"
	t.Log(messageCode.String(), string(messageCode))
	t.Log(messageCode.For("en-us"))
}
