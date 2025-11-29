// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package userconfigs

import "github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"

type UserFeatureCode = string

const (
	UserFeatureCodePlan                UserFeatureCode = "plan"
	UserFeatureCodeServerTCP           UserFeatureCode = "server.tcp"
	UserFeatureCodeServerTCPPort       UserFeatureCode = "server.tcp.port"
	UserFeatureCodeServerUDP           UserFeatureCode = "server.udp"
	UserFeatureCodeServerUDPPort       UserFeatureCode = "server.udp.port"
	UserFeatureCodeServerAccessLog     UserFeatureCode = "server.accessLog"
	UserFeatureCodeServerViewAccessLog UserFeatureCode = "server.viewAccessLog"
	UserFeatureCodeServerScript        UserFeatureCode = "server.script"
	UserFeatureCodeServerWAF           UserFeatureCode = "server.waf"
	UserFeatureCodeServerOptimization  UserFeatureCode = "server.optimization"
	UserFeatureCodeServerUAM           UserFeatureCode = "server.uam"
	UserFeatureCodeServerWebP          UserFeatureCode = "server.webp"
	UserFeatureCodeServerCC            UserFeatureCode = "server.cc"
	UserFeatureCodeServerACME          UserFeatureCode = "server.acme"
	UserFeatureCodeServerAuth          UserFeatureCode = "server.auth"
	UserFeatureCodeServerWebsocket     UserFeatureCode = "server.websocket"
	UserFeatureCodeServerHTTP3         UserFeatureCode = "server.http3"
	UserFeatureCodeServerReferers      UserFeatureCode = "server.referers"
	UserFeatureCodeServerUserAgent     UserFeatureCode = "server.userAgent"
	UserFeatureCodeServerRequestLimit  UserFeatureCode = "server.requestLimit"
	UserFeatureCodeServerCompression   UserFeatureCode = "server.compression"
	UserFeatureCodeServerRewriteRules  UserFeatureCode = "server.rewriteRules"
	UserFeatureCodeServerHostRedirects UserFeatureCode = "server.hostRedirects"
	UserFeatureCodeServerHTTPHeaders   UserFeatureCode = "server.httpHeaders"
	UserFeatureCodeServerPages         UserFeatureCode = "server.pages"
)

type UserFeature struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	SupportPlan bool   `json:"supportPlan"`
}

func (this *UserFeature) ToPB() *pb.UserFeature {
	return &pb.UserFeature{
		Name:        this.Name,
		Code:        this.Code,
		Description: this.Description,
		SupportPlan: this.SupportPlan,
	}
}

func FindAllUserFeatures() []*UserFeature {
	return []*UserFeature{
		{Name: "Access log", Code: UserFeatureCodeServerAccessLog, Description: "allow writing access log", SupportPlan: true},
		{Name: "View access log", Code: UserFeatureCodeServerViewAccessLog, Description: "allow viewing access log", SupportPlan: true},
		{Name: "TCP load balancer", Code: UserFeatureCodeServerTCP, Description: "allow creating TCP/TLS servers", SupportPlan: false},
		{Name: "Custom TCP ports", Code: UserFeatureCodeServerTCPPort, Description: "allow custom TCP ports", SupportPlan: true},
		{Name: "UDP load balancer", Code: UserFeatureCodeServerUDP, Description: "allow creating UDP servers", SupportPlan: false},
		{Name: "Custom UDP ports", Code: UserFeatureCodeServerUDPPort, Description: "allow custom UDP ports", SupportPlan: true},
		{Name: "ACME certificates", Code: UserFeatureCodeServerACME, Description: "allow requesting free ACME certificates", SupportPlan: false},
		{Name: "WAF", Code: UserFeatureCodeServerWAF, Description: "enable web application firewall", SupportPlan: true},
		{Name: "Edge script", Code: UserFeatureCodeServerScript, Description: "run edge scripts", SupportPlan: true},
		{Name: "5s shield", Code: UserFeatureCodeServerUAM, Description: "enable challenge page", SupportPlan: true},
		{Name: "CC protection", Code: UserFeatureCodeServerCC, Description: "enable CC protection", SupportPlan: true},
		{Name: "WebP", Code: UserFeatureCodeServerWebP, Description: "auto convert WebP", SupportPlan: true},
		{Name: "Optimization", Code: UserFeatureCodeServerOptimization, Description: "page optimization", SupportPlan: true},
		{Name: "Auth", Code: UserFeatureCodeServerAuth, Description: "enable request auth", SupportPlan: true},
		{Name: "WebSocket", Code: UserFeatureCodeServerWebsocket, Description: "enable WebSocket", SupportPlan: true},
		{Name: "Anti-leech", Code: UserFeatureCodeServerReferers, Description: "referer whitelist/blacklist", SupportPlan: true},
		{Name: "UA control", Code: UserFeatureCodeServerUserAgent, Description: "user-agent whitelist/blacklist", SupportPlan: true},
		{Name: "HTTP/3", Code: UserFeatureCodeServerHTTP3, Description: "enable HTTP/3", SupportPlan: true},
		{Name: "Request limit", Code: UserFeatureCodeServerRequestLimit, Description: "limit concurrency/bandwidth", SupportPlan: true},
		{Name: "Compression", Code: UserFeatureCodeServerCompression, Description: "enable compression", SupportPlan: true},
		{Name: "URL redirect", Code: UserFeatureCodeServerHostRedirects, Description: "configure redirects", SupportPlan: true},
		{Name: "Rewrite rules", Code: UserFeatureCodeServerRewriteRules, Description: "custom rewrite rules", SupportPlan: true},
		{Name: "HTTP headers", Code: UserFeatureCodeServerHTTPHeaders, Description: "manage request/response headers", SupportPlan: true},
		{Name: "Custom pages", Code: UserFeatureCodeServerPages, Description: "custom error pages", SupportPlan: true},
		{Name: "Plan", Code: UserFeatureCodePlan, Description: "buy/manage plans", SupportPlan: false},
	}
}

func FindUserFeature(code string) *UserFeature {
	for _, feature := range FindAllUserFeatures() {
		if feature.Code == code {
			return feature
		}
	}
	return nil
}

func CheckUserFeature(featureCode string) bool {
	return FindUserFeature(featureCode) != nil
}
