package posts

import (
	"github.com/TeaOSLab/EdgeAPI/internal/infra/db"
	"gorm.io/gorm"
)

const (
	PostCategoryStateEnabled  = 1
	PostCategoryStateDisabled = 0
)

type PostCategoryDAO struct {
	db    *gorm.DB
	Table string
}

func NewPostCategoryDAO() (*PostCategoryDAO, error) {
	conn, err := db.DB()
	if err != nil {
		return nil, err
	}
	return &PostCategoryDAO{db: conn, Table: "edgePostCategories"}, nil
}

var SharedPostCategoryDAO *PostCategoryDAO

func init() {
	dao, err := NewPostCategoryDAO()
	if err == nil {
		SharedPostCategoryDAO = dao
	}
}

func (dao *PostCategoryDAO) EnablePostCategory(categoryId int64) error {
	return dao.db.Model(&PostCategory{}).
		Where("id=?", categoryId).
		Update("state", PostCategoryStateEnabled).Error
}

func (dao *PostCategoryDAO) DisablePostCategory(categoryId int64) error {
	return dao.db.Model(&PostCategory{}).
		Where("id=?", categoryId).
		Update("state", PostCategoryStateDisabled).Error
}

func (dao *PostCategoryDAO) FindEnabledPostCategory(categoryId int64) (*PostCategory, error) {
	var cat PostCategory
	err := dao.db.Where("id=? AND state=?", categoryId, PostCategoryStateEnabled).First(&cat).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cat, nil
}

func (dao *PostCategoryDAO) FindPostCategoryName(categoryId int64) (string, error) {
	var name string
	err := dao.db.Model(&PostCategory{}).
		Where("id=?", categoryId).
		Select("name").
		Scan(&name).Error
	return name, err
}
