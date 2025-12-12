package dns

import (
	"errors"
	"time"

	dbutils "github.com/TeaOSLab/EdgeAPI/internal/db/utils"
	"github.com/TeaOSLab/EdgeAPI/internal/infra/db"
	"github.com/iwind/TeaGo/dbs"
	"gorm.io/gorm"
)

const (
	DNSProviderStateEnabled  = 1 // 已启用
	DNSProviderStateDisabled = 0 // 已禁用
)

type DNSProviderDAO struct {
	db *gorm.DB
}

func NewDNSProviderDAO() (*DNSProviderDAO, error) {
	conn, err := db.DB()
	if err != nil {
		return nil, err
	}
	return &DNSProviderDAO{db: conn}, nil
}

func (dao *DNSProviderDAO) useDB(_ *dbs.Tx) *gorm.DB {
	return dao.db
}

var SharedDNSProviderDAO *DNSProviderDAO

func init() {
	dao, err := NewDNSProviderDAO()
	if err == nil {
		SharedDNSProviderDAO = dao
	}
}

// EnableDNSProvider 启用条目
func (dao *DNSProviderDAO) EnableDNSProvider(tx *dbs.Tx, id int64) error {
	return dao.useDB(tx).
		Model(&DNSProvider{}).
		Where("id = ?", id).
		Update("state", DNSProviderStateEnabled).
		Error
}

// DisableDNSProvider 禁用条目
func (dao *DNSProviderDAO) DisableDNSProvider(tx *dbs.Tx, id int64) error {
	return dao.useDB(tx).
		Model(&DNSProvider{}).
		Where("id = ?", id).
		Update("state", DNSProviderStateDisabled).
		Error
}

// FindEnabledDNSProvider 查找启用中的条目
func (dao *DNSProviderDAO) FindEnabledDNSProvider(tx *dbs.Tx, id int64) (*DNSProvider, error) {
	var provider DNSProvider
	err := dao.useDB(tx).
		Where("id = ? AND state = ?", id, DNSProviderStateEnabled).
		First(&provider).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &provider, err
}

// CreateDNSProvider 创建服务商
func (dao *DNSProviderDAO) CreateDNSProvider(tx *dbs.Tx, adminId int64, userId int64, providerType string, name string, apiParamsJSON []byte, minTTL int32) (int64, error) {
	provider := &DNSProvider{
		AdminId:   adminId,
		UserId:    userId,
		Type:      providerType,
		Name:      name,
		ApiParams: apiParamsJSON,
		State:     DNSProviderStateEnabled,
		CreatedAt: time.Now().Unix(),
	}
	if minTTL >= 0 {
		provider.MinTTL = minTTL
	}
	if len(apiParamsJSON) == 0 {
		provider.ApiParams = nil
	}

	if err := dao.useDB(tx).Create(provider).Error; err != nil {
		return 0, err
	}
	return provider.Id, nil
}

// UpdateDNSProvider 修改服务商
func (dao *DNSProviderDAO) UpdateDNSProvider(tx *dbs.Tx, dnsProviderId int64, name string, apiParamsJSON []byte, minTTL int32) error {
	if dnsProviderId <= 0 {
		return errors.New("invalid dnsProviderId")
	}

	updates := map[string]any{
		"name": name,
	}
	if len(apiParamsJSON) > 0 {
		updates["apiParams"] = apiParamsJSON
	}
	if minTTL >= 0 {
		updates["minTTL"] = minTTL
	}

	return dao.useDB(tx).
		Model(&DNSProvider{}).
		Where("id = ?", dnsProviderId).
		Updates(updates).
		Error
}

// CountAllEnabledDNSProviders 计算服务商数量
func (dao *DNSProviderDAO) CountAllEnabledDNSProviders(tx *dbs.Tx, adminId int64, userId int64, keyword string, domain string, providerType string) (int64, error) {
	dbQuery := dao.useDB(tx).Model(&DNSProvider{}).Where("state = ?", DNSProviderStateEnabled)
	if userId > 0 {
		dbQuery = dbQuery.Where("userId = ?", userId)
	}
	if len(keyword) > 0 {
		dbQuery = dbQuery.Where("name LIKE ?", dbutils.QuoteLike(keyword))
	}
	if len(domain) > 0 {
		subQuery := dao.useDB(tx).
			Model(&DNSDomain{}).
			Select("providerId").
			Where("state = ? AND name = ?", DNSDomainStateEnabled, domain)
		dbQuery = dbQuery.Where("id IN (?)", subQuery)
	}
	if len(providerType) > 0 {
		dbQuery = dbQuery.Where("type = ?", providerType)
	}

	var count int64
	err := dbQuery.Count(&count).Error
	return count, err
}

// ListEnabledDNSProviders 列出单页服务商
func (dao *DNSProviderDAO) ListEnabledDNSProviders(tx *dbs.Tx, adminId int64, userId int64, keyword string, domain string, providerType string, offset int64, size int64) (result []*DNSProvider, err error) {
	dbQuery := dao.useDB(tx).Model(&DNSProvider{}).Where("state = ?", DNSProviderStateEnabled)
	if userId > 0 {
		dbQuery = dbQuery.Where("userId = ?", userId)
	}
	if len(keyword) > 0 {
		dbQuery = dbQuery.Where("name LIKE ?", dbutils.QuoteLike(keyword))
	}
	if len(domain) > 0 {
		subQuery := dao.useDB(tx).
			Model(&DNSDomain{}).
			Select("providerId").
			Where("state = ? AND name = ?", DNSDomainStateEnabled, domain)
		dbQuery = dbQuery.Where("id IN (?)", subQuery)
	}
	if len(providerType) > 0 {
		dbQuery = dbQuery.Where("type = ?", providerType)
	}

	err = dbQuery.
		Order("id DESC").
		Offset(int(offset)).
		Limit(int(size)).
		Find(&result).
		Error
	return
}

// FindAllEnabledDNSProviders 列出所有服务商
func (dao *DNSProviderDAO) FindAllEnabledDNSProviders(tx *dbs.Tx, adminId int64, userId int64) (result []*DNSProvider, err error) {
	dbQuery := dao.useDB(tx).Model(&DNSProvider{}).Where("state = ?", DNSProviderStateEnabled)
	if userId > 0 {
		dbQuery = dbQuery.Where("userId = ?", userId)
	}
	err = dbQuery.
		Order("id DESC").
		Find(&result).
		Error
	return
}

// FindAllEnabledDNSProvidersWithType 查询某个类型下的所有服务商
func (dao *DNSProviderDAO) FindAllEnabledDNSProvidersWithType(tx *dbs.Tx, providerType string) (result []*DNSProvider, err error) {
	err = dao.useDB(tx).
		Where("state = ? AND type = ?", DNSProviderStateEnabled, providerType).
		Order("id DESC").
		Find(&result).
		Error
	return
}

// UpdateProviderDataUpdatedTime 更新数据更新时间
func (dao *DNSProviderDAO) UpdateProviderDataUpdatedTime(tx *dbs.Tx, providerId int64) error {
	return dao.useDB(tx).
		Model(&DNSProvider{}).
		Where("id = ?", providerId).
		Update("dataUpdatedAt", time.Now().Unix()).
		Error
}
