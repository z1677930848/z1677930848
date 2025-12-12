package dns

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/TeaOSLab/EdgeAPI/internal/dnsclients/dnstypes"
	"github.com/TeaOSLab/EdgeAPI/internal/infra/db"
	"github.com/TeaOSLab/EdgeAPI/internal/utils"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/types"
	"gorm.io/gorm"
)

const (
	DNSDomainStateEnabled  = 1 // 已启用
	DNSDomainStateDisabled = 0 // 已禁用
)

const DNSDomainTableName = "edgeDNSDomains"

// DNSDomain represents a managed DNS domain.
type DNSDomain struct {
	Id            int64           `gorm:"column:id;primaryKey"`
	AdminId       int64           `gorm:"column:adminId"`
	UserId        int64           `gorm:"column:userId"`
	ProviderId    int64           `gorm:"column:providerId"`
	IsOn          bool            `gorm:"column:isOn"`
	Name          string          `gorm:"column:name"`
	CreatedAt     int64           `gorm:"column:createdAt"`
	DataUpdatedAt int64           `gorm:"column:dataUpdatedAt"`
	DataError     string          `gorm:"column:dataError"`
	Data          json.RawMessage `gorm:"column:data"`
	Records       json.RawMessage `gorm:"column:records"`
	Routes        json.RawMessage `gorm:"column:routes"`
	IsUp          bool            `gorm:"column:isUp"`
	State         int             `gorm:"column:state"`
	IsDeleted     bool            `gorm:"column:isDeleted"`
	DnsName       string          `gorm:"column:dnsName"`
	Accounts      json.RawMessage `gorm:"column:accounts"`
	DnsInfo       json.RawMessage `gorm:"column:dnsInfo"`
	RegionId      int64           `gorm:"column:regionId"`
	UpdatedAt     time.Time       `gorm:"column:updatedAt"`
}

func (DNSDomain) TableName() string { return DNSDomainTableName }

type DNSDomainDAO struct {
	db    *gorm.DB
	Table string
}

func NewDNSDomainDAO() (*DNSDomainDAO, error) {
	conn, err := db.DB()
	if err != nil {
		return nil, err
	}
	return &DNSDomainDAO{db: conn, Table: DNSDomainTableName}, nil
}

func (dao *DNSDomainDAO) useDB(_ *dbs.Tx) *gorm.DB {
	return dao.db
}

// EnableDNSDomain 启用条目
func (dao *DNSDomainDAO) EnableDNSDomain(tx *dbs.Tx, id int64) error {
	return dao.useDB(tx).Model(&DNSDomain{}).
		Where("id = ?", id).
		Update("state", DNSDomainStateEnabled).Error
}

// DisableDNSDomain 禁用条目
func (dao *DNSDomainDAO) DisableDNSDomain(tx *dbs.Tx, id int64) error {
	return dao.useDB(tx).Model(&DNSDomain{}).
		Where("id = ?", id).
		Update("state", DNSDomainStateDisabled).Error
}

// FindEnabledDNSDomain 查找启用中的条目
func (dao *DNSDomainDAO) FindEnabledDNSDomain(tx *dbs.Tx, domainId int64, cacheMap *utils.CacheMap) (*DNSDomain, error) {
	cacheKey := DNSDomainTableName + ":record:" + types.String(domainId)
	if cacheMap != nil {
		if v, _ := cacheMap.Get(cacheKey); v != nil {
			return v.(*DNSDomain), nil
		}
	}

	var domain DNSDomain
	err := dao.useDB(tx).Where("id = ? AND state = ?", domainId, DNSDomainStateEnabled).First(&domain).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if cacheMap != nil {
		cacheMap.Put(cacheKey, &domain)
	}
	return &domain, nil
}

// FindDNSDomainName 根据主键查找名称
func (dao *DNSDomainDAO) FindDNSDomainName(tx *dbs.Tx, id int64) (string, error) {
	var name string
	err := dao.useDB(tx).Model(&DNSDomain{}).Where("id = ?", id).Select("name").Scan(&name).Error
	return name, err
}

// CreateDomain 创建域名
func (dao *DNSDomainDAO) CreateDomain(tx *dbs.Tx, adminId int64, userId int64, providerId int64, name string) (int64, error) {
	domain := &DNSDomain{
		ProviderId: providerId,
		AdminId:    adminId,
		UserId:     userId,
		Name:       name,
		State:      DNSDomainStateEnabled,
		IsOn:       true,
		IsUp:       true,
		CreatedAt:  time.Now().Unix(),
	}
	if err := dao.useDB(tx).Create(domain).Error; err != nil {
		return 0, err
	}
	return domain.Id, nil
}

// UpdateDomain 修改域名
func (dao *DNSDomainDAO) UpdateDomain(tx *dbs.Tx, domainId int64, name string, isOn bool) error {
	if domainId <= 0 {
		return errors.New("invalid domainId")
	}
	return dao.useDB(tx).Model(&DNSDomain{}).
		Where("id = ?", domainId).
		Updates(map[string]any{
			"name": name,
			"isOn": isOn,
		}).Error
}

// FindAllEnabledDomainsWithProviderId 查询一个服务商下面的所有域名
func (dao *DNSDomainDAO) FindAllEnabledDomainsWithProviderId(tx *dbs.Tx, providerId int64) (result []*DNSDomain, err error) {
	err = dao.useDB(tx).
		Where("providerId = ? AND state = ?", providerId, DNSDomainStateEnabled).
		Order("id ASC").
		Find(&result).
		Error
	return
}

// ListDomains 列出单页域名
func (dao *DNSDomainDAO) ListDomains(tx *dbs.Tx, providerId int64, isDeleted bool, isUp bool, offset int64, size int64) (result []*DNSDomain, err error) {
	err = dao.useDB(tx).
		Where("providerId = ? AND state = ? AND isDeleted = ? AND isUp = ?", providerId, DNSDomainStateEnabled, isDeleted, isUp).
		Order("id ASC").
		Offset(int(offset)).
		Limit(int(size)).
		Find(&result).
		Error
	return
}

// CountAllEnabledDomainsWithProviderId 计算某个服务商下的域名数量
func (dao *DNSDomainDAO) CountAllEnabledDomainsWithProviderId(tx *dbs.Tx, providerId int64, isDeleted bool, isUp bool) (int64, error) {
	var count int64
	err := dao.useDB(tx).
		Model(&DNSDomain{}).
		Where("providerId = ? AND state = ? AND isDeleted = ? AND isUp = ?", providerId, DNSDomainStateEnabled, isDeleted, isUp).
		Count(&count).
		Error
	return count, err
}

// UpdateDomainData 更新域名数据
func (dao *DNSDomainDAO) UpdateDomainData(tx *dbs.Tx, domainId int64, data string) error {
	if domainId <= 0 {
		return errors.New("invalid domainId")
	}
	return dao.useDB(tx).Model(&DNSDomain{}).
		Where("id = ?", domainId).
		Update("data", data).
		Error
}

// UpdateDomainRecords 更新域名解析记录
func (dao *DNSDomainDAO) UpdateDomainRecords(tx *dbs.Tx, domainId int64, recordsJSON []byte) error {
	if domainId <= 0 {
		return errors.New("invalid domainId")
	}
	return dao.useDB(tx).Model(&DNSDomain{}).
		Where("id = ?", domainId).
		Updates(map[string]any{
			"records":       recordsJSON,
			"dataUpdatedAt": time.Now().Unix(),
		}).Error
}

// UpdateDomainRoutes 更新线路
func (dao *DNSDomainDAO) UpdateDomainRoutes(tx *dbs.Tx, domainId int64, routesJSON []byte) error {
	if domainId <= 0 {
		return errors.New("invalid domainId")
	}
	return dao.useDB(tx).Model(&DNSDomain{}).
		Where("id = ?", domainId).
		Updates(map[string]any{
			"routes":        routesJSON,
			"dataUpdatedAt": time.Now().Unix(),
		}).Error
}

// FindDomainRoutes 查找域名线路
func (dao *DNSDomainDAO) FindDomainRoutes(tx *dbs.Tx, domainId int64) ([]*dnstypes.Route, error) {
	var routesData []byte
	err := dao.useDB(tx).
		Model(&DNSDomain{}).
		Where("id = ?", domainId).
		Select("routes").
		Scan(&routesData).
		Error
	if err != nil {
		return nil, err
	}
	if len(routesData) == 0 || string(routesData) == "null" {
		return nil, nil
	}
	var result []*dnstypes.Route
	if err := json.Unmarshal(routesData, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// FindDomainRouteName 查找线路名称
func (dao *DNSDomainDAO) FindDomainRouteName(tx *dbs.Tx, domainId int64, routeCode string) (string, error) {
	routes, err := dao.FindDomainRoutes(tx, domainId)
	if err != nil {
		return "", err
	}
	for _, route := range routes {
		if route.Code == routeCode {
			return route.Name, nil
		}
	}
	return "", nil
}

// ExistAvailableDomains 判断是否有域名可选
func (dao *DNSDomainDAO) ExistAvailableDomains(tx *dbs.Tx) (bool, error) {
	subQuery := dao.useDB(tx).
		Model(&DNSProvider{}).
		Select("id").
		Where("state = ?", DNSProviderStateEnabled)

	var count int64
	err := dao.useDB(tx).
		Model(&DNSDomain{}).
		Where("state = ? AND isOn = ? AND providerId IN (?)", DNSDomainStateEnabled, true, subQuery).
		Count(&count).
		Error
	return count > 0, err
}

// ExistDomainRecord 检查域名解析记录是否存在
func (dao *DNSDomainDAO) ExistDomainRecord(tx *dbs.Tx, domainId int64, recordName string, recordType string, recordRoute string, recordValue string) (bool, error) {
	recordType = strings.ToUpper(recordType)

	query := map[string]any{
		"name": recordName,
		"type": recordType,
	}
	if len(recordRoute) > 0 {
		query["route"] = recordRoute
	}
	if len(recordValue) > 0 {
		query["value"] = recordValue

		// CNAME兼容点（.）符号
		if recordType == "CNAME" && !strings.HasSuffix(recordValue, ".") {
			b, err := dao.ExistDomainRecord(tx, domainId, recordName, recordType, recordRoute, recordValue+".")
			if err != nil {
				return false, err
			}
			if b {
				return true, nil
			}
		}
	}
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return false, err
	}

	var count int64
	err = dao.useDB(tx).
		Model(&DNSDomain{}).
		Where("id = ?", domainId).
		Where("JSON_CONTAINS(records, ?)", queryJSON).
		Count(&count).
		Error
	return count > 0, err
}

// FindEnabledDomainWithName 根据名称查找某个域名
func (dao *DNSDomainDAO) FindEnabledDomainWithName(tx *dbs.Tx, providerId int64, domainName string) (*DNSDomain, error) {
	var domain DNSDomain
	err := dao.useDB(tx).
		Where("state = ? AND providerId = ? AND name = ?", DNSDomainStateEnabled, providerId, domainName).
		First(&domain).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &domain, err
}

// UpdateDomainIsUp 设置是否在线
func (dao *DNSDomainDAO) UpdateDomainIsUp(tx *dbs.Tx, domainId int64, isUp bool) error {
	return dao.useDB(tx).
		Model(&DNSDomain{}).
		Where("id = ?", domainId).
		Update("isUp", isUp).
		Error
}

// UpdateDomainIsDeleted 设置域名为删除
func (dao *DNSDomainDAO) UpdateDomainIsDeleted(tx *dbs.Tx, domainId int64, isDeleted bool) error {
	return dao.useDB(tx).
		Model(&DNSDomain{}).
		Where("id = ?", domainId).
		Update("isDeleted", isDeleted).
		Error
}

// DecodeRoutes 解析线路配置
func (d *DNSDomain) DecodeRoutes() ([]*dnstypes.Route, error) {
	if len(d.Routes) == 0 || string(d.Routes) == "null" {
		return nil, nil
	}
	var result []*dnstypes.Route
	if err := json.Unmarshal(d.Routes, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ContainsRouteCode 检查线路代码是否存在
func (d *DNSDomain) ContainsRouteCode(route string) (bool, error) {
	routes, err := d.DecodeRoutes()
	if err != nil {
		return false, err
	}
	for _, r := range routes {
		if r.Code == route {
			return true, nil
		}
	}
	return false, nil
}

// DecodeRecords 解析记录配置
func (d *DNSDomain) DecodeRecords() ([]*dnstypes.Record, error) {
	if len(d.Records) == 0 || string(d.Records) == "null" {
		return nil, nil
	}
	var result []*dnstypes.Record
	if err := json.Unmarshal(d.Records, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Init global dao (compatible with legacy pattern).
var SharedDNSDomainDAO *DNSDomainDAO

func init() {
	dao, err := NewDNSDomainDAO()
	if err == nil {
		SharedDNSDomainDAO = dao
	}
}
