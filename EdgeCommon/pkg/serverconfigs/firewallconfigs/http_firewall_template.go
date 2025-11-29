package firewallconfigs

import "github.com/iwind/TeaGo/maps"

type HTTPFirewallRuleConnector = string

const (
	HTTPFirewallRuleConnectorAnd = "and"
	HTTPFirewallRuleConnectorOr  = "or"
)

// HTTPFirewallTemplate returns a tiny builtin policy
func HTTPFirewallTemplate() *HTTPFirewallPolicy {
	policy := &HTTPFirewallPolicy{
		IsOn:     true,
		Inbound:  &HTTPFirewallInboundConfig{},
		Outbound: &HTTPFirewallOutboundConfig{},
	}

	// group: XSS
	xss := &HTTPFirewallRuleGroup{
		IsOn:        true,
		Name:        "XSS",
		Code:        "xss",
		Description: "Block cross-site scripting",
		IsTemplate:  true,
	}
	set := &HTTPFirewallRuleSet{
		IsOn:      true,
		Name:      "XSS detection",
		Code:      "1010",
		Connector: HTTPFirewallRuleConnectorOr,
		Actions: []*HTTPFirewallActionConfig{
			{Code: HTTPFirewallActionPage, Options: maps.Map{"status": 403, "body": ""}},
		},
	}
	set.AddRule(&HTTPFirewallRule{IsOn: true, Param: "${requestAll}", Operator: HTTPFirewallRuleOperatorContainsXSS, Value: "", IsCaseInsensitive: false})
	xss.AddRuleSet(set)
	policy.Inbound.Groups = append(policy.Inbound.Groups, xss)

	// group: file upload
	upload := &HTTPFirewallRuleGroup{
		IsOn:        false,
		Name:        "Upload",
		Code:        "upload",
		Description: "Block executable file upload",
		IsTemplate:  true,
	}
	uploadSet := &HTTPFirewallRuleSet{
		IsOn:      true,
		Name:      "Dangerous extensions",
		Code:      "2001",
		Connector: HTTPFirewallRuleConnectorOr,
		Actions: []*HTTPFirewallActionConfig{
			{Code: HTTPFirewallActionPage, Options: maps.Map{"status": 403, "body": ""}},
		},
	}
	uploadSet.AddRule(&HTTPFirewallRule{IsOn: true, Param: "${requestUpload.ext}", Operator: HTTPFirewallRuleOperatorMatch, Value: `\.(php|jsp|aspx|asp|exe|asa|rb|py)$`, IsCaseInsensitive: true})
	upload.AddRuleSet(uploadSet)
	policy.Inbound.Groups = append(policy.Inbound.Groups, upload)

	return policy
}
