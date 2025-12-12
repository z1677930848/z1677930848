package dns

const DNSTaskTableName = "edgeDNSTasks"

// DNSTask DNS更新任务
type DNSTask struct {
	Id         int64  `gorm:"column:id;primaryKey"` // ID
	ClusterId  int64  `gorm:"column:clusterId"`     // 集群ID
	ServerId   int64  `gorm:"column:serverId"`      // 服务ID
	NodeId     int64  `gorm:"column:nodeId"`        // 节点ID
	DomainId   int64  `gorm:"column:domainId"`      // 域名ID
	RecordName string `gorm:"column:recordName"`    // 记录名
	Type       string `gorm:"column:type"`          // 任务类型
	UpdatedAt  int64  `gorm:"column:updatedAt"`     // 更新时间
	IsDone     bool   `gorm:"column:isDone"`        // 是否已完成
	IsOk       bool   `gorm:"column:isOk"`          // 是否成功
	Error      string `gorm:"column:error"`         // 错误信息
	Version    int64  `gorm:"column:version"`       // 版本
	CountFails int    `gorm:"column:countFails"`    // 尝试失败次数
}

func (DNSTask) TableName() string { return DNSTaskTableName }
