// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package langs

type MessageCode string

func (this MessageCode) For(langCode LangCode, args ...any) string {
	return Message(langCode, this, args...)
}

func (this MessageCode) String() string {
	return string(this)
}
