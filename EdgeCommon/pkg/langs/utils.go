// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package langs

import (
	"net/http"
	"strings"

	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/iwind/TeaGo/actions"
)

// Message 璇诲彇娑堟伅
// Read message
func Message(langCode LangCode, messageCode MessageCode, args ...any) string {
	return defaultManager.GetMessage(langCode, messageCode, args...)
}

func DefaultMessage(messageCode MessageCode, args ...any) string {
	return defaultManager.GetMessage("en-us", messageCode, args...)
}

func ParseLangFromRequest(req *http.Request) (langCode string) {
	// parse language from cookie
	const cookieName = "edgelang"
	cookie, _ := req.Cookie(cookieName)
	if cookie != nil && len(cookie.Value) > 0 && defaultManager.HasLang(cookie.Value) {
		return cookie.Value
	}

	// parse language from 'Accept-Language'
	var acceptLanguage = req.Header.Get("Accept-Language")
	if len(acceptLanguage) > 0 {
		var pieces = strings.Split(acceptLanguage, ",")
		for _, lang := range pieces {
			var index = strings.Index(lang, ";")
			if index >= 0 {
				lang = lang[:index]
			}

			var match = defaultManager.MatchLang(lang)
			if len(match) > 0 {
				return match
			}
		}
	}

	return defaultManager.DefaultLang()
}

func ParseLangFromAction(action actions.ActionWrapper) (langCode string) {
	return ParseLangFromRequest(action.Object().Request)
}

// Format 鏍煎紡鍖栧彉閲?
// Format string that contains message variables, such as ${lang.MESSAGE_CODE}
//
// 鏆傛椂涓嶆敮鎸佸彉閲忎腑鍔犲弬鏁?
func Format(langCode LangCode, varString string) string {
	return configutils.ParseVariables(varString, func(varName string) (value string) {
		if !strings.HasPrefix(varName, varPrefix) {
			return "${" + varName + "}" // keep origin variable
		}
		return Message(langCode, MessageCode(varName[len(varPrefix):]))
	})
}

// Load 鍔犺浇娑堟伅瀹氫箟
// Load message definitions from map
func Load(langCode LangCode, messageMap map[MessageCode]string) {
	lang, ok := defaultManager.GetLang(langCode)
	if !ok {
		lang = defaultManager.AddLang(langCode)
	}
	for messageCode, messageText := range messageMap {
		lang.Set(messageCode, messageText)
	}
}
