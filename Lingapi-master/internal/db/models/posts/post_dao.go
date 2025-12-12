package posts

import (
	"github.com/TeaOSLab/EdgeAPI/internal/infra/db"
	"gorm.io/gorm"
)

const (
	PostStateEnabled  = 1
	PostStateDisabled = 0
)

type PostDAO struct {
	db    *gorm.DB
	Table string
}

func NewPostDAO() (*PostDAO, error) {
	conn, err := db.DB()
	if err != nil {
		return nil, err
	}
	return &PostDAO{db: conn, Table: "edgePosts"}, nil
}

var SharedPostDAO *PostDAO

func init() {
	dao, err := NewPostDAO()
	if err == nil {
		SharedPostDAO = dao
	}
}

func (dao *PostDAO) EnablePost(postId int64) error {
	return dao.db.Model(&Post{}).
		Where("id=?", postId).
		Update("state", PostStateEnabled).Error
}

func (dao *PostDAO) DisablePost(postId int64) error {
	return dao.db.Model(&Post{}).
		Where("id=?", postId).
		Update("state", PostStateDisabled).Error
}

func (dao *PostDAO) FindEnabledPost(postId int64) (*Post, error) {
	var p Post
	err := dao.db.Where("id=? AND state=?", postId, PostStateEnabled).First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
