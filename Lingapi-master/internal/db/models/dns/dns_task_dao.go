package dns

import (
	"errors"
	"time"

	"github.com/TeaOSLab/EdgeAPI/internal/infra/db"
	"github.com/iwind/TeaGo/dbs"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DNSTaskType = string

const (
	DNSTaskTypeClusterChange       DNSTaskType = "clusterChange"       // 集群节点、服务发生变化
	DNSTaskTypeClusterNodesChange  DNSTaskType = "clusterNodesChange"  // 集群中节点发生变化
	DNSTaskTypeClusterRemoveDomain DNSTaskType = "clusterRemoveDomain" // 从集群中移除域名
	DNSTaskTypeNodeChange          DNSTaskType = "nodeChange"
	DNSTaskTypeServerChange        DNSTaskType = "serverChange"
	DNSTaskTypeDomainChange        DNSTaskType = "domainChange"
)

var DNSTasksNotifier = make(chan bool, 2)

type DNSTaskDAO struct {
	db *gorm.DB
}

func NewDNSTaskDAO() (*DNSTaskDAO, error) {
	conn, err := db.DB()
	if err != nil {
		return nil, err
	}
	return &DNSTaskDAO{db: conn}, nil
}

func (dao *DNSTaskDAO) useDB(_ *dbs.Tx) *gorm.DB {
	return dao.db
}

var SharedDNSTaskDAO *DNSTaskDAO

func init() {
	dao, err := NewDNSTaskDAO()
	if err == nil {
		SharedDNSTaskDAO = dao
	}
}

// CreateDNSTask 生成任务
func (dao *DNSTaskDAO) CreateDNSTask(tx *dbs.Tx, clusterId int64, serverId int64, nodeId int64, domainId int64, recordName string, taskType string) error {
	if clusterId <= 0 && serverId <= 0 && nodeId <= 0 && domainId <= 0 {
		return nil
	}

	now := time.Now().Unix()
	version := time.Now().UnixNano()

	task := &DNSTask{
		ClusterId:  clusterId,
		ServerId:   serverId,
		NodeId:     nodeId,
		DomainId:   domainId,
		RecordName: recordName,
		Type:       taskType,
		UpdatedAt:  now,
		IsDone:     false,
		IsOk:       false,
		Error:      "",
		Version:    version,
		CountFails: 0,
	}

	err := dao.useDB(tx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "clusterId"},
				{Name: "serverId"},
				{Name: "nodeId"},
				{Name: "domainId"},
				{Name: "recordName"},
				{Name: "type"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"updatedAt":  now,
				"isDone":     false,
				"isOk":       false,
				"error":      "",
				"version":    version,
				"countFails": 0,
			}),
		}).
		Create(task).Error
	if err != nil {
		return err
	}

	select {
	case DNSTasksNotifier <- true:
	default:
	}
	return nil
}

// CreateClusterTask 生成集群变更任务
func (dao *DNSTaskDAO) CreateClusterTask(tx *dbs.Tx, clusterId int64, taskType DNSTaskType) error {
	return dao.CreateDNSTask(tx, clusterId, 0, 0, 0, "", taskType)
}

// CreateClusterRemoveTask 生成集群删除域名任务
func (dao *DNSTaskDAO) CreateClusterRemoveTask(tx *dbs.Tx, clusterId int64, domainId int64, recordName string) error {
	return dao.CreateDNSTask(tx, clusterId, 0, 0, domainId, recordName, DNSTaskTypeClusterRemoveDomain)
}

// CreateNodeTask 生成节点任务
func (dao *DNSTaskDAO) CreateNodeTask(tx *dbs.Tx, clusterId int64, nodeId int64, taskType DNSTaskType) error {
	return dao.CreateDNSTask(tx, clusterId, 0, nodeId, 0, "", taskType)
}

// CreateServerTask 生成服务任务
func (dao *DNSTaskDAO) CreateServerTask(tx *dbs.Tx, clusterId int64, serverId int64, taskType DNSTaskType) error {
	return dao.CreateDNSTask(tx, clusterId, serverId, 0, 0, "", taskType)
}

// CreateDomainTask 生成域名更新任务
func (dao *DNSTaskDAO) CreateDomainTask(tx *dbs.Tx, domainId int64, taskType DNSTaskType) error {
	return dao.CreateDNSTask(tx, 0, 0, 0, domainId, "", taskType)
}

// FindAllDoingTasks 查找所有正在执行的任务
func (dao *DNSTaskDAO) FindAllDoingTasks(tx *dbs.Tx) (result []*DNSTask, err error) {
	err = dao.useDB(tx).
		Where("(isDone = 0 OR (isDone = 1 AND isOk = 0 AND countFails < 3))").
		Order("version ASC").
		Order("id ASC").
		Find(&result).
		Error
	return
}

// FindAllDoingOrErrorTasks 查找正在执行的和错误的任务
func (dao *DNSTaskDAO) FindAllDoingOrErrorTasks(tx *dbs.Tx, nodeClusterId int64) (result []*DNSTask, err error) {
	dbQuery := dao.useDB(tx).Model(&DNSTask{}).
		Where("(isDone = 0 OR (isDone = 1 AND isOk = 0))")
	if nodeClusterId > 0 {
		dbQuery = dbQuery.Where("clusterId = ?", nodeClusterId)
	}
	err = dbQuery.
		Order("updatedAt ASC").
		Order("version ASC").
		Order("id ASC").
		Find(&result).
		Error
	return
}

// ExistDoingTasks 检查是否有正在执行的任务
func (dao *DNSTaskDAO) ExistDoingTasks(tx *dbs.Tx) (bool, error) {
	var count int64
	err := dao.useDB(tx).
		Model(&DNSTask{}).
		Where("isDone = 0").
		Count(&count).
		Error
	return count > 0, err
}

// ExistErrorTasks 检查是否有错误的任务
func (dao *DNSTaskDAO) ExistErrorTasks(tx *dbs.Tx) (bool, error) {
	var count int64
	err := dao.useDB(tx).
		Model(&DNSTask{}).
		Where("isDone = 1 AND isOk = 0").
		Count(&count).
		Error
	return count > 0, err
}

// DeleteDNSTask 删除任务
func (dao *DNSTaskDAO) DeleteDNSTask(tx *dbs.Tx, taskId int64) error {
	return dao.useDB(tx).
		Where("id = ?", taskId).
		Delete(&DNSTask{}).
		Error
}

// DeleteAllDNSTasks 删除所有任务
func (dao *DNSTaskDAO) DeleteAllDNSTasks(tx *dbs.Tx) error {
	return dao.useDB(tx).
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&DNSTask{}).
		Error
}

// UpdateDNSTaskError 设置任务错误
func (dao *DNSTaskDAO) UpdateDNSTaskError(tx *dbs.Tx, taskId int64, errMsg string) error {
	if taskId <= 0 {
		return errors.New("invalid taskId")
	}
	return dao.useDB(tx).
		Model(&DNSTask{}).
		Where("id = ?", taskId).
		Updates(map[string]any{
			"isDone":     true,
			"isOk":       false,
			"error":      errMsg,
			"countFails": gorm.Expr("countFails+?", 1),
		}).Error
}

// UpdateDNSTaskDone 设置任务完成
func (dao *DNSTaskDAO) UpdateDNSTaskDone(tx *dbs.Tx, taskId int64, taskVersion int64) error {
	if taskId <= 0 {
		return errors.New("invalid taskId")
	}

	var currentVersion int64
	err := dao.useDB(tx).
		Model(&DNSTask{}).
		Where("id = ?", taskId).
		Select("version").
		Scan(&currentVersion).Error
	if err != nil {
		return err
	}

	if taskVersion > 0 && currentVersion > 0 && currentVersion != taskVersion {
		return nil
	}

	return dao.useDB(tx).
		Model(&DNSTask{}).
		Where("id = ?", taskId).
		Updates(map[string]any{
			"isDone":     true,
			"isOk":       true,
			"countFails": 0,
			"error":      "",
		}).Error
}

// GenerateVersion 生成最新的版本号
func (dao *DNSTaskDAO) GenerateVersion() int64 {
	return time.Now().UnixNano()
}

// UpdateClusterDNSTasksDone 设置所有集群任务完成
func (dao *DNSTaskDAO) UpdateClusterDNSTasksDone(tx *dbs.Tx, clusterId int64, maxVersion int64) error {
	if clusterId <= 0 || maxVersion <= 0 {
		return nil
	}

	return dao.useDB(tx).
		Model(&DNSTask{}).
		Where("clusterId = ? AND isOk = ? AND version <= ?", clusterId, false, maxVersion).
		Updates(map[string]any{
			"isDone":     true,
			"isOk":       true,
			"error":      "",
			"countFails": 0,
		}).Error
}

// DeleteDNSTasksWithClusterId 删除集群相关任务
func (dao *DNSTaskDAO) DeleteDNSTasksWithClusterId(tx *dbs.Tx, clusterId int64) error {
	if clusterId <= 0 {
		return nil
	}
	return dao.useDB(tx).
		Where("clusterId = ?", clusterId).
		Delete(&DNSTask{}).
		Error
}
