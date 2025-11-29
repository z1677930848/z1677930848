package firewallconfigs

type HTTPFirewallRuleOperator = string

type HTTPFirewallRuleCaseInsensitive = string

const (
	HTTPFirewallRuleOperatorGt                           HTTPFirewallRuleOperator = "gt"
	HTTPFirewallRuleOperatorGte                          HTTPFirewallRuleOperator = "gte"
	HTTPFirewallRuleOperatorLt                           HTTPFirewallRuleOperator = "lt"
	HTTPFirewallRuleOperatorLte                          HTTPFirewallRuleOperator = "lte"
	HTTPFirewallRuleOperatorEq                           HTTPFirewallRuleOperator = "eq"
	HTTPFirewallRuleOperatorNeq                          HTTPFirewallRuleOperator = "neq"
	HTTPFirewallRuleOperatorEqString                     HTTPFirewallRuleOperator = "eq string"
	HTTPFirewallRuleOperatorNeqString                    HTTPFirewallRuleOperator = "neq string"
	HTTPFirewallRuleOperatorMatch                        HTTPFirewallRuleOperator = "match"
	HTTPFirewallRuleOperatorNotMatch                     HTTPFirewallRuleOperator = "not match"
	HTTPFirewallRuleOperatorWildcardMatch                HTTPFirewallRuleOperator = "wildcard match"
	HTTPFirewallRuleOperatorWildcardNotMatch             HTTPFirewallRuleOperator = "wildcard not match"
	HTTPFirewallRuleOperatorContains                     HTTPFirewallRuleOperator = "contains"
	HTTPFirewallRuleOperatorNotContains                  HTTPFirewallRuleOperator = "not contains"
	HTTPFirewallRuleOperatorContainsAnyWord              HTTPFirewallRuleOperator = "contains any word"
	HTTPFirewallRuleOperatorContainsAllWords             HTTPFirewallRuleOperator = "contains all words"
	HTTPFirewallRuleOperatorNotContainsAnyWord           HTTPFirewallRuleOperator = "not contains any word"
	HTTPFirewallRuleOperatorPrefix                       HTTPFirewallRuleOperator = "prefix"
	HTTPFirewallRuleOperatorSuffix                       HTTPFirewallRuleOperator = "suffix"
	HTTPFirewallRuleOperatorContainsAny                  HTTPFirewallRuleOperator = "contains any"
	HTTPFirewallRuleOperatorContainsAll                  HTTPFirewallRuleOperator = "contains all"
	HTTPFirewallRuleOperatorContainsSQLInjection         HTTPFirewallRuleOperator = "contains sql injection"
	HTTPFirewallRuleOperatorContainsSQLInjectionStrictly HTTPFirewallRuleOperator = "contains sql injection strictly"
	HTTPFirewallRuleOperatorContainsXSS                  HTTPFirewallRuleOperator = "contains xss"
	HTTPFirewallRuleOperatorContainsXSSStrictly          HTTPFirewallRuleOperator = "contains xss strictly"
	HTTPFirewallRuleOperatorHasKey                       HTTPFirewallRuleOperator = "has key"
	HTTPFirewallRuleOperatorVersionGt                    HTTPFirewallRuleOperator = "version gt"
	HTTPFirewallRuleOperatorVersionLt                    HTTPFirewallRuleOperator = "version lt"
	HTTPFirewallRuleOperatorVersionRange                 HTTPFirewallRuleOperator = "version range"
	HTTPFirewallRuleOperatorContainsBinary               HTTPFirewallRuleOperator = "contains binary"
	HTTPFirewallRuleOperatorNotContainsBinary            HTTPFirewallRuleOperator = "not contains binary"
	HTTPFirewallRuleOperatorEqIP                         HTTPFirewallRuleOperator = "eq ip"
	HTTPFirewallRuleOperatorInIPList                     HTTPFirewallRuleOperator = "in ip list"
	HTTPFirewallRuleOperatorGtIP                         HTTPFirewallRuleOperator = "gt ip"
	HTTPFirewallRuleOperatorGteIP                        HTTPFirewallRuleOperator = "gte ip"
	HTTPFirewallRuleOperatorLtIP                         HTTPFirewallRuleOperator = "lt ip"
	HTTPFirewallRuleOperatorLteIP                        HTTPFirewallRuleOperator = "lte ip"
	HTTPFirewallRuleOperatorIPRange                      HTTPFirewallRuleOperator = "ip range"
	HTTPFirewallRuleOperatorNotIPRange                   HTTPFirewallRuleOperator = "not ip range"
	HTTPFirewallRuleOperatorIPMod10                      HTTPFirewallRuleOperator = "ip mod 10"
	HTTPFirewallRuleOperatorIPMod100                     HTTPFirewallRuleOperator = "ip mod 100"
	HTTPFirewallRuleOperatorIPMod                        HTTPFirewallRuleOperator = "ip mod"

	HTTPFirewallRuleCaseInsensitiveNone = "none"
	HTTPFirewallRuleCaseInsensitiveYes  = "yes"
	HTTPFirewallRuleCaseInsensitiveNo   = "no"
)

type RuleOperatorDefinition struct {
	Name            string                          `json:"name"`
	Code            string                          `json:"code"`
	Description     string                          `json:"description"`
	CaseInsensitive HTTPFirewallRuleCaseInsensitive `json:"caseInsensitive"`
	DataType        string                          `json:"dataType"`
}

var AllRuleOperators = []*RuleOperatorDefinition{
	{Name: "regexp match", Code: HTTPFirewallRuleOperatorMatch, Description: "regular expression match", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveYes, DataType: "regexp"},
	{Name: "regexp not match", Code: HTTPFirewallRuleOperatorNotMatch, Description: "regular expression not match", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveYes, DataType: "regexp"},
	{Name: "wildcard match", Code: HTTPFirewallRuleOperatorWildcardMatch, Description: "match using * wildcard", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveYes, DataType: "wildcard"},
	{Name: "wildcard not match", Code: HTTPFirewallRuleOperatorWildcardNotMatch, Description: "not match using wildcard", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveYes, DataType: "wildcard"},
	{Name: "string equal", Code: HTTPFirewallRuleOperatorEqString, Description: "string equals", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNo, DataType: "string"},
	{Name: "string not equal", Code: HTTPFirewallRuleOperatorNeqString, Description: "string not equals", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNo, DataType: "string"},
	{Name: "contains", Code: HTTPFirewallRuleOperatorContains, Description: "contains substring", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNo, DataType: "string"},
	{Name: "not contains", Code: HTTPFirewallRuleOperatorNotContains, Description: "not contains substring", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNo, DataType: "string"},
	{Name: "contains any", Code: HTTPFirewallRuleOperatorContainsAny, Description: "contains any string per line", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNo, DataType: "strings"},
	{Name: "contains all", Code: HTTPFirewallRuleOperatorContainsAll, Description: "contains all strings", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNo, DataType: "strings"},
	{Name: "prefix", Code: HTTPFirewallRuleOperatorPrefix, Description: "has prefix", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNo, DataType: "string"},
	{Name: "suffix", Code: HTTPFirewallRuleOperatorSuffix, Description: "has suffix", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNo, DataType: "string"},
	{Name: "contains any word", Code: HTTPFirewallRuleOperatorContainsAnyWord, Description: "contains word list", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNo, DataType: "strings"},
	{Name: "contains all words", Code: HTTPFirewallRuleOperatorContainsAllWords, Description: "contains all words", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNo, DataType: "strings"},
	{Name: "not contains any word", Code: HTTPFirewallRuleOperatorNotContainsAnyWord, Description: "no words matched", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNo, DataType: "strings"},
	{Name: "contains SQL injection", Code: HTTPFirewallRuleOperatorContainsSQLInjection, Description: "detect SQL injection", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "none"},
	{Name: "contains SQL injection strictly", Code: HTTPFirewallRuleOperatorContainsSQLInjectionStrictly, Description: "strict SQL injection detect", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "none"},
	{Name: "contains XSS", Code: HTTPFirewallRuleOperatorContainsXSS, Description: "detect XSS", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "none"},
	{Name: "contains XSS strictly", Code: HTTPFirewallRuleOperatorContainsXSSStrictly, Description: "strict XSS detect", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "none"},
	{Name: "contains binary", Code: HTTPFirewallRuleOperatorContainsBinary, Description: "contains binary data", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNo, DataType: "string"},
	{Name: "not contains binary", Code: HTTPFirewallRuleOperatorNotContainsBinary, Description: "not contains binary", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNo, DataType: "string"},
	{Name: "number greater", Code: HTTPFirewallRuleOperatorGt, Description: "numeric greater than", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "number"},
	{Name: "number greater or equal", Code: HTTPFirewallRuleOperatorGte, Description: "numeric >=", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "number"},
	{Name: "number less", Code: HTTPFirewallRuleOperatorLt, Description: "numeric <", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "number"},
	{Name: "number less or equal", Code: HTTPFirewallRuleOperatorLte, Description: "numeric <=", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "number"},
	{Name: "number equal", Code: HTTPFirewallRuleOperatorEq, Description: "numeric ==", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "number"},
	{Name: "number not equal", Code: HTTPFirewallRuleOperatorNeq, Description: "numeric !=", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "number"},
	{Name: "has key", Code: HTTPFirewallRuleOperatorHasKey, Description: "slice/map contains key", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNo, DataType: "string"},
	{Name: "version greater", Code: HTTPFirewallRuleOperatorVersionGt, Description: "greater version", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "version"},
	{Name: "version less", Code: HTTPFirewallRuleOperatorVersionLt, Description: "less version", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "version"},
	{Name: "version range", Code: HTTPFirewallRuleOperatorVersionRange, Description: "version range start,end", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "versionRange"},
	{Name: "ip equal", Code: HTTPFirewallRuleOperatorEqIP, Description: "single IP equals", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "ip"},
	{Name: "ip in list", Code: HTTPFirewallRuleOperatorInIPList, Description: "in IP list", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "ips"},
	{Name: "ip greater", Code: HTTPFirewallRuleOperatorGtIP, Description: "IP greater", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "ip"},
	{Name: "ip greater or equal", Code: HTTPFirewallRuleOperatorGteIP, Description: "IP >=", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "ip"},
	{Name: "ip less", Code: HTTPFirewallRuleOperatorLtIP, Description: "IP <", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "ip"},
	{Name: "ip less or equal", Code: HTTPFirewallRuleOperatorLteIP, Description: "IP <=", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "ip"},
	{Name: "ip range", Code: HTTPFirewallRuleOperatorIPRange, Description: "ip in range or CIDR", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "ips"},
	{Name: "ip not range", Code: HTTPFirewallRuleOperatorNotIPRange, Description: "ip not in range", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "ips"},
	{Name: "ip mod 10", Code: HTTPFirewallRuleOperatorIPMod10, Description: "ip value mod 10", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "number"},
	{Name: "ip mod 100", Code: HTTPFirewallRuleOperatorIPMod100, Description: "ip value mod 100", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "number"},
	{Name: "ip mod", Code: HTTPFirewallRuleOperatorIPMod, Description: "ip mod divisor,remainder", CaseInsensitive: HTTPFirewallRuleCaseInsensitiveNone, DataType: "number"},
}

func FindRuleOperatorName(code string) string {
	for _, operator := range AllRuleOperators {
		if operator.Code == code {
			return operator.Name
		}
	}
	return ""
}
