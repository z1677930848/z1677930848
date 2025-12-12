package posts

type PostCategory struct {
	Id     int64  `gorm:"column:id;primaryKey"`
	Name   string `gorm:"column:name"`
	IsOn   bool   `gorm:"column:isOn"`
	Code   string `gorm:"column:code"`
	Order  uint32 `gorm:"column:order"`
	State  int    `gorm:"column:state"`
}

func (PostCategory) TableName() string {
	return "edgePostCategories"
}
