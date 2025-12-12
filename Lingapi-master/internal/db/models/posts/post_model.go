package posts

type Post struct {
	Id          int64  `gorm:"column:id;primaryKey"`
	CategoryId  int64  `gorm:"column:categoryId"`
	Type        string `gorm:"column:type"`
	Url         string `gorm:"column:url"`
	Subject     string `gorm:"column:subject"`
	Body        string `gorm:"column:body"`
	CreatedAt   int64  `gorm:"column:createdAt"`
	IsPublished bool   `gorm:"column:isPublished"`
	PublishedAt int64  `gorm:"column:publishedAt"`
	ProductCode string `gorm:"column:productCode"`
	State       int    `gorm:"column:state"`
}

func (Post) TableName() string {
	return "edgePosts"
}
