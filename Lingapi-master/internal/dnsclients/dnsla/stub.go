//go:build plus

package dnsla

// 仅用于占位编译，真实实现请在商业版中提供。

type DomainListResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    DomainListData    `json:"data"`
}

type DomainListData struct {
	List []Domain `json:"list"`
}

type Domain struct {
	DomainId int64  `json:"domain_id"`
	Domain   string `json:"domain"`
}

type RecordListResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    RecordListData   `json:"data"`
}

type RecordListData struct {
	List []Record `json:"list"`
}

type Record struct {
	RecordId int64  `json:"record_id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Route    string `json:"route"`
	TTL      int    `json:"ttl"`
}

type AllLineListResponse struct {
	Data AllLineListResponseData `json:"data"`
}

type AllLineListResponseData struct {
	Children []AllLineListResponseChild `json:"children"`
}

type AllLineListResponseChild struct {
	Name     string                        `json:"name"`
	Route    string                        `json:"route"`
	Children []AllLineListResponseChild    `json:"children"`
}

type RecordCreateResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		RecordId int64 `json:"record_id"`
	} `json:"data"`
}

type RecordUpdateResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type RecordDeleteResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type DomainResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
