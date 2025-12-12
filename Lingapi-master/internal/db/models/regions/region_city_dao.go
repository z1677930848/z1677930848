package regions

import (
	"encoding/json"
	"sort"

	"github.com/TeaOSLab/EdgeAPI/internal/infra/db"
	"github.com/TeaOSLab/EdgeAPI/internal/utils"
	"github.com/TeaOSLab/EdgeAPI/internal/utils/numberutils"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/types"
	"gorm.io/gorm"
)

const (
	RegionCityStateEnabled  = 1
	RegionCityStateDisabled = 0
)

type RegionCityDAO struct {
	db    *gorm.DB
	Table string
}

func NewRegionCityDAO() (*RegionCityDAO, error) {
	conn, err := db.DB()
	if err != nil {
		return nil, err
	}
	return &RegionCityDAO{db: conn, Table: "edgeRegionCities"}, nil
}

var SharedRegionCityDAO *RegionCityDAO

func init() {
	dao, err := NewRegionCityDAO()
	if err == nil {
		SharedRegionCityDAO = dao
	}
}

func (dao *RegionCityDAO) EnableRegionCity(id uint32) error {
	return dao.db.Model(&RegionCity{}).
		Where("valueId = ?", id).
		Update("state", RegionCityStateEnabled).Error
}

func (dao *RegionCityDAO) DisableRegionCity(id uint32) error {
	return dao.db.Model(&RegionCity{}).
		Where("valueId = ?", id).
		Update("state", RegionCityStateDisabled).Error
}

func (dao *RegionCityDAO) FindEnabledRegionCity(_ *dbs.Tx, id int64) (*RegionCity, error) {
	var city RegionCity
	err := dao.db.Where("valueId=? AND state=?", id, RegionCityStateEnabled).First(&city).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &city, nil
}

func (dao *RegionCityDAO) FindRegionCityName(id uint32) (string, error) {
	var name string
	err := dao.db.Model(&RegionCity{}).Where("valueId=?", id).Select("name").Scan(&name).Error
	return name, err
}

func (dao *RegionCityDAO) FindCityWithDataId(dataId string) (int64, error) {
	var valueId int64
	err := dao.db.Model(&RegionCity{}).Where("dataId=?", dataId).Select("valueId").Scan(&valueId).Error
	return valueId, err
}

func (dao *RegionCityDAO) CreateCity(provinceId int64, name string, dataId string) (int64, error) {
	codesJSON, err := json.Marshal([]string{name})
	if err != nil {
		return 0, err
	}
	city := &RegionCity{
		ProvinceId: provinceId,
		Name:       name,
		DataId:     dataId,
		State:      RegionCityStateEnabled,
		Codes:      codesJSON,
	}
	if err := dao.db.Create(city).Error; err != nil {
		return 0, err
	}
	// set valueId = id
	_ = dao.db.Model(city).Update("valueId", city.Id).Error
	return city.Id, nil
}

func (dao *RegionCityDAO) FindCityIdWithName(_ *dbs.Tx, provinceId int64, cityName string) (int64, error) {
	var valueId int64
	err := dao.db.Model(&RegionCity{}).
		Where("provinceId=? AND (name=? OR customName=?)", provinceId, cityName, cityName).
		Select("valueId").
		Scan(&valueId).Error
	return valueId, err
}

func (dao *RegionCityDAO) FindAllEnabledCities(_ *dbs.Tx) ([]*RegionCity, error) {
	var result []*RegionCity
	err := dao.db.Where("state=?", RegionCityStateEnabled).Find(&result).Error
	return result, err
}

func (dao *RegionCityDAO) FindAllEnabledCitiesWithProvinceId(_ *dbs.Tx, provinceId int64) ([]*RegionCity, error) {
	var result []*RegionCity
	err := dao.db.Where("provinceId=? AND state=?", provinceId, RegionCityStateEnabled).Find(&result).Error
	return result, err
}

func (dao *RegionCityDAO) UpdateCityCustom(_ *dbs.Tx, cityId int64, customName string, customCodes []string) error {
	if customCodes == nil {
		customCodes = []string{}
	}
	customCodesJSON, err := json.Marshal(customCodes)
	if err != nil {
		return err
	}
	return dao.db.Model(&RegionCity{}).
		Where("valueId=?", cityId).
		Updates(map[string]any{
			"customName":  customName,
			"customCodes": customCodesJSON,
		}).Error
}

// Similarity helper
func (dao *RegionCityDAO) FindSimilarCities(cities []*RegionCity, cityName string, size int) (result []*RegionCity) {
	if len(cities) == 0 {
		return
	}
	type scoreCity struct {
		score float32
		city  *RegionCity
	}
	var scored []scoreCity
	for _, city := range cities {
		var similarityList []float32
		for _, code := range city.AllCodes() {
			sim := utils.Similar(cityName, code)
			if sim > 0 {
				similarityList = append(similarityList, sim)
			}
		}
		if len(similarityList) > 0 {
			scored = append(scored, scoreCity{
				score: numberutils.Max(similarityList...),
				city:  city,
			})
		}
	}
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})
	if len(scored) > size {
		scored = scored[:size]
	}
	for _, s := range scored {
		result = append(result, s.city)
	}
	return
}

// Optional helper to match old signature
func (dao *RegionCityDAO) FindEnabledRegionCityName(id int64) (string, error) {
	name, err := dao.FindRegionCityName(uint32(id))
	return name, err
}

// Converts string to int for previous calls expecting int param
func parseInt(v int64) int {
	return types.Int(v)
}
