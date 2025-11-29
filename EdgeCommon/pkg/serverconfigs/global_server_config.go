// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package serverconfigs

type LnRequestSchedulingMethod = string

const (
	LnRequestSchedulingMethodRandom     LnRequestSchedulingMethod = "random"
	LnRequestSchedulingMethodURLMapping LnRequestSchedulingMethod = "urlMapping"
)

const (
	DefaultTCPPortRangeMin = 10000
	DefaultTCPPortRangeMax = 40000
)

func NewGlobalServerConfig() *GlobalServerConfig {
	var config = &GlobalServerConfig{}

	config.HTTPAll.SupportsLowVersionHTTP = true
	config.HTTPAll.EnableServerAddrVariable = false
	config.HTTPAll.LnRequestSchedulingMethod = LnRequestSchedulingMethodURLMapping

	config.HTTPAccessLog.IsOn = true
	config.HTTPAccessLog.EnableRequestHeaders = true
	config.HTTPAccessLog.EnableResponseHeaders = true
	config.HTTPAccessLog.EnableCookies = true
	config.HTTPAccessLog.EnableServerNotFound = true

	config.Log.RecordServerError = false

	config.Performance.AutoWriteTimeout = false
	config.Performance.AutoReadTimeout = false
	config.Stat.Upload.MaxCities = 32
	config.Stat.Upload.MaxProviders = 32
	config.Stat.Upload.MaxSystems = 64
	config.Stat.Upload.MaxBrowsers = 64

	return config
}

// GlobalServerConfig 鍏ㄥ眬鐨勬湇鍔￠厤缃?
type GlobalServerConfig struct {
	HTTPAll struct {
		MatchDomainStrictly  bool                  `yaml:"matchDomainStrictly" json:"matchDomainStrictly"`   // 鏄惁涓ユ牸鍖归厤鍩熷悕
		AllowMismatchDomains []string              `yaml:"allowMismatchDomains" json:"allowMismatchDomains"` // 鍏佽鐨勪笉鍖归厤鐨勫煙鍚?
		AllowNodeIP          bool                  `yaml:"allowNodeIP" json:"allowNodeIP"`                   // 鍏佽IP鐩存帴璁块棶
		NodeIPShowPage       bool                  `yaml:"nodeIPShowPage" json:"nodeIPShowPage"`             // 璁块棶IP鍦板潃鏄惁鏄剧ず椤甸潰
		NodeIPPageHTML       string                `yaml:"nodeIPPageHTML" json:"nodeIPPageHTML"`             // 璁块棶IP鍦板潃椤甸潰鍐呭
		DefaultDomain        string                `yaml:"defaultDomain" json:"defaultDomain"`               // 榛樿鐨勫煙鍚?
		DomainMismatchAction *DomainMismatchAction `yaml:"domainMismatchAction" json:"domainMismatchAction"` // 涓嶅尮閰嶆椂閲囧彇鐨勫姩浣?

		SupportsLowVersionHTTP    bool                      `yaml:"supportsLowVersionHTTP" json:"supportsLowVersionHTTP"`       // 鏄惁鍚敤浣庣増鏈琀TTP
		MatchCertFromAllServers   bool                      `yaml:"matchCertFromAllServers" json:"matchCertFromAllServers"`     // 浠庢墍鏈夋湇鍔′腑鍖归厤璇佷功锛堜笉瑕佽交鏄撳紑鍚紒锛?
		ForceLnRequest            bool                      `yaml:"forceLnRequest" json:"forceLnRequest"`                       // 寮哄埗浠嶭n璇锋眰鍐呭
		LnRequestSchedulingMethod LnRequestSchedulingMethod `yaml:"lnRequestSchedulingMethod" json:"lnRequestSchedulingMethod"` // Ln璇锋眰璋冨害鏂规硶
		ServerName                string                    `yaml:"serverName" json:"serverName"`                               // Server鍚嶇О
		EnableServerAddrVariable  bool                      `yaml:"enableServerAddrVariable" json:"enableServerAddrVariable"`   // 鏄惁鏀寔${serverAddr}鍙橀噺
		XFFMaxAddresses           int                       `yaml:"xffMaxAddresses" json:"xffMaxAddresses"`                     // XFF涓渶澶氱殑鍦板潃鏁?

		DomainAuditingIsOn   bool   `yaml:"domainAuditingIsOn" json:"domainAuditingIsOn"`     // 鍩熷悕鏄惁闇€瑕佸鏍?
		DomainAuditingPrompt string `yaml:"domainAuditingPrompt" json:"domainAuditingPrompt"` // 鍩熷悕瀹℃牳鐨勬彁绀?

		RequestOriginsWithEncodings bool `yaml:"requestOriginsWithEncodings" json:"requestOriginsWithEncodings"` // 浣跨敤浣跨敤鍘嬬缉缂栫爜鍥炴簮
	} `yaml:"httpAll" json:"httpAll"` // HTTP缁熶竴閰嶇疆

	TCPAll struct {
		PortRangeMin int   `yaml:"portRangeMin" json:"portRangeMin"` // 鏈€灏忕鍙?
		PortRangeMax int   `yaml:"portRangeMax" json:"portRangeMax"` // 鏈€澶х鍙?
		DenyPorts    []int `yaml:"denyPorts" json:"denyPorts"`       // 绂佹浣跨敤鐨勭鍙?
	} `yaml:"tcpAll" json:"tcpAll"`

	HTTPAccessLog struct {
		IsOn                     bool `yaml:"isOn" json:"isOn"`                                         // 鏄惁鍚敤姝ゅ姛鑳?
		EnableRequestHeaders     bool `yaml:"enableRequestHeaders" json:"enableRequestHeaders"`         // 璁板綍璇锋眰Header
		CommonRequestHeadersOnly bool `yaml:"commonRequestHeadersOnly" json:"commonRequestHeadersOnly"` // 鍙繚鐣欓€氱敤Header
		EnableResponseHeaders    bool `yaml:"enableResponseHeaders" json:"enableResponseHeaders"`       // 璁板綍鍝嶅簲Header
		EnableCookies            bool `yaml:"enableCookies" json:"enableCookies"`                       // 璁板綍Cookie
		EnableServerNotFound     bool `yaml:"enableServerNotFound" json:"enableServerNotFound"`         // 璁板綍鏈嶅姟鎵句笉鍒扮殑鏃ュ織
	} `yaml:"httpAccessLog" json:"httpAccessLog"` // 璁块棶鏃ュ織閰嶇疆

	Stat struct {
		Upload struct {
			MaxCities    int16 `yaml:"maxCities" json:"maxCities"`       // 鏈€澶у尯鍩熸暟閲?
			MaxProviders int16 `yaml:"maxProviders" json:"maxProviders"` // 鏈€澶ц繍钀ュ晢鏁伴噺
			MaxSystems   int16 `yaml:"maxSystems" json:"maxSystems"`     // 鏈€澶ф搷浣滅郴缁熸暟閲?
			MaxBrowsers  int16 `yaml:"maxBrowsers" json:"maxBrowsers"`   // 鏈€澶ф祻瑙堝櫒鏁伴噺
		} `yaml:"upload" json:"upload"` // 涓婁紶鐩稿叧璁剧疆
	} `yaml:"stat" json:"stat"` // 缁熻鐩稿叧閰嶇疆

	Performance struct {
		Debug            bool `yaml:"debug" json:"debug"`                       // Debug妯″紡
		AutoWriteTimeout bool `yaml:"autoWriteTimeout" json:"autoWriteTimeout"` // 鏄惁鑷姩鍐欒秴鏃?
		AutoReadTimeout  bool `yaml:"autoReadTimeout" json:"autoReadTimeout"`   // 鏄惁鑷姩璇昏秴鏃?
	} `yaml:"performance" json:"performance"` // 鎬ц兘

	Log struct {
		RecordServerError bool `yaml:"recordServerError" json:"recordServerError"` // 璁板綍鏈嶅姟閿欒鍒拌繍琛屾棩蹇?
	} `yaml:"log" json:"log"` // 杩愯鏃ュ織閰嶇疆
}

func (this *GlobalServerConfig) Init() error {
	// 鏈壘鍒板煙鍚嶆椂鐨勫姩浣?
	if this.HTTPAll.DomainMismatchAction != nil {
		err := this.HTTPAll.DomainMismatchAction.Init()
		if err != nil {
			return err
		}
	}

	return nil
}
