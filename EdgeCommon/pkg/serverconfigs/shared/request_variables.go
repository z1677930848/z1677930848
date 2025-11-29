// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package shared

import "github.com/iwind/TeaGo/maps"

// DefaultRequestVariables returns built-in request variables
func DefaultRequestVariables() []maps.Map {
	return []maps.Map{
		{"code": "${edgeVersion}", "name": "edge version", "description": ""},
		{"code": "${remoteAddr}", "name": "client address", "description": "read from X-Forwarded-For/X-Real-IP/RemoteAddr"},
		{"code": "${rawRemoteAddr}", "name": "raw client address", "description": "direct client IP"},
		{"code": "${remotePort}", "name": "client port", "description": ""},
		{"code": "${remoteUser}", "name": "client username", "description": ""},
		{"code": "${requestURI}", "name": "request URI", "description": "e.g. /hello?name=lily"},
		{"code": "${requestPath}", "name": "request path", "description": "without query string"},
		{"code": "${requestURL}", "name": "request URL", "description": "e.g. https://example.com/hello?name=lily"},
		{"code": "${requestLength}", "name": "content length", "description": ""},
		{"code": "${requestMethod}", "name": "method", "description": "GET/POST/..."},
		{"code": "${requestFilename}", "name": "filename", "description": ""},
		{"code": "${requestPathExtension}", "name": "path extension", "description": "with dot"},
		{"code": "${requestPathLowerExtension}", "name": "path extension lower", "description": "with dot"},
		{"code": "${scheme}", "name": "scheme", "description": "http or https"},
		{"code": "${proto}", "name": "protocol", "description": "e.g. HTTP/1.1"},
		{"code": "${timeISO8601}", "name": "ISO8601 time", "description": "e.g. 2018-07-16T23:52:24+08:00"},
		{"code": "${timeLocal}", "name": "local time", "description": "e.g. 17/Jul/2018:09:52:24 +0800"},
		{"code": "${msec}", "name": "unix msec", "description": "seconds.milliseconds"},
		{"code": "${timestamp}", "name": "unix timestamp", "description": "seconds"},
		{"code": "${host}", "name": "host", "description": ""},
		{"code": "${cname}", "name": "site CNAME", "description": ""},
		{"code": "${serverName}", "name": "server name", "description": ""},
		{"code": "${serverPort}", "name": "server port", "description": ""},
		{"code": "${referer}", "name": "referer", "description": ""},
		{"code": "${referer.host}", "name": "referer host", "description": ""},
		{"code": "${userAgent}", "name": "user agent", "description": ""},
		{"code": "${contentType}", "name": "content type", "description": ""},
		{"code": "${cookies}", "name": "cookies", "description": "raw cookie string"},
		{"code": "${cookie.NAME}", "name": "cookie", "description": "NAME is cookie key"},
		{"code": "${isArgs}", "name": "arg flag", "description": "? if query string exists"},
		{"code": "${args}", "name": "args", "description": "full query string"},
		{"code": "${arg.NAME}", "name": "arg", "description": "single query parameter"},
		{"code": "${headers}", "name": "headers", "description": "all headers"},
		{"code": "${header.NAME}", "name": "header", "description": "single header"},
		{"code": "${geo.country.name}", "name": "country name", "description": ""},
		{"code": "${geo.country.id}", "name": "country id", "description": ""},
		{"code": "${geo.province.name}", "name": "province name", "description": "China only"},
		{"code": "${geo.province.id}", "name": "province id", "description": "China only"},
		{"code": "${geo.city.name}", "name": "city name", "description": "China only"},
		{"code": "${geo.city.id}", "name": "city id", "description": "China only"},
		{"code": "${isp.name}", "name": "ISP name", "description": ""},
		{"code": "${isp.id}", "name": "ISP id", "description": ""},
		{"code": "${browser.os.name}", "name": "OS name", "description": ""},
		{"code": "${browser.os.version}", "name": "OS version", "description": ""},
		{"code": "${browser.name}", "name": "browser name", "description": ""},
		{"code": "${browser.version}", "name": "browser version", "description": ""},
		{"code": "${browser.isMobile}", "name": "is mobile", "description": "1 for mobile, 0 otherwise"},
	}
}
