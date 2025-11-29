// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package dnsconfigs

type RecordType = string

const (
	RecordTypeA     RecordType = "A"
	RecordTypeCNAME RecordType = "CNAME"
	RecordTypeAAAA  RecordType = "AAAA"
	RecordTypeNS    RecordType = "NS"
	RecordTypeMX    RecordType = "MX"
	RecordTypeSRV   RecordType = "SRV"
	RecordTypeTXT   RecordType = "TXT"
	RecordTypeCAA   RecordType = "CAA"
	RecordTypeSOA   RecordType = "SOA"
)

type RecordTypeDefinition struct {
	Type        RecordType `json:"type"`
	Description string     `json:"description"`
	CanDefine   bool       `json:"canDefine"` // whether users can create this type
}

func FindAllRecordTypeDefinitions() []*RecordTypeDefinition {
	return []*RecordTypeDefinition{
		{Type: RecordTypeA, Description: "map hostname to an IPv4 address", CanDefine: true},
		{Type: RecordTypeCNAME, Description: "map hostname to another hostname", CanDefine: true},
		{Type: RecordTypeAAAA, Description: "map hostname to an IPv6 address", CanDefine: true},
		{Type: RecordTypeNS, Description: "delegate sub-domain to another DNS server", CanDefine: false},
		{Type: RecordTypeSOA, Description: "start of authority", CanDefine: false},
		{Type: RecordTypeMX, Description: "mail exchange record", CanDefine: true},
		{Type: RecordTypeSRV, Description: "service location record", CanDefine: true},
		{Type: RecordTypeTXT, Description: "text/SPF/verification record", CanDefine: true},
		{Type: RecordTypeCAA, Description: "CA authorization", CanDefine: true},
	}
}

func FindAllUserRecordTypeDefinitions() []*RecordTypeDefinition {
	var result []*RecordTypeDefinition
	for _, r := range FindAllRecordTypeDefinitions() {
		if r.CanDefine {
			result = append(result, r)
		}
	}
	return result
}
