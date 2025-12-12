package dns

import "gorm.io/datatypes"

const DNSProviderTableName = "edgeDNSProviders"

// DNSProvider DNS服务商
type DNSProvider struct {
	Id            int64          `gorm:"column:id;primaryKey"` // ID
	Name          string         `gorm:"column:name"`          // 名称
	AdminId       int64          `gorm:"column:adminId"`       // 管理员ID
	UserId        int64          `gorm:"column:userId"`        // 用户ID
	Type          string         `gorm:"column:type"`          // 供应商类型
	ApiParams     datatypes.JSON `gorm:"column:apiParams"`     // API参数
	CreatedAt     int64          `gorm:"column:createdAt"`     // 创建时间
	State         int            `gorm:"column:state"`         // 状态
	DataUpdatedAt int64          `gorm:"column:dataUpdatedAt"` // 数据同步时间
	MinTTL        int32          `gorm:"column:minTTL"`        // 最小TTL
}

func (DNSProvider) TableName() string { return DNSProviderTableName }
