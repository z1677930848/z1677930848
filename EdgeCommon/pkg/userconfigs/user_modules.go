// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package userconfigs

type UserModule = string

const (
	UserModuleCDN      UserModule = "cdn"
	UserModuleAntiDDoS UserModule = "antiDDoS"
	UserModuleNS       UserModule = "ns"
)

var DefaultUserModules = []UserModule{UserModuleCDN}
