package regions

type RegionCity struct {
	Id          int64  `gorm:"column:id;primaryKey"`
	ValueId     int64  `gorm:"column:valueId"`
	ProvinceId  int64  `gorm:"column:provinceId"`
	Name        string `gorm:"column:name"`
	Codes       []byte `gorm:"column:codes"`
	CustomName  string `gorm:"column:customName"`
	CustomCodes []byte `gorm:"column:customCodes"`
	State       int    `gorm:"column:state"`
	DataId      string `gorm:"column:dataId"`
	IsOn        bool   `gorm:"column:isOn"`
	IsDeleted   bool   `gorm:"column:isDeleted"`
}

func (RegionCity) TableName() string {
	return "edgeRegionCities"
}
