package model

import (
	"gin-api/pkg/app"
	"github.com/jinzhu/gorm"
)

type ArticleSwagger struct {
	List  []*Article
	Pager *app.Pager
}

type Article struct {
	Id           uint64 `json:"id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Introduction string `json:"introduction"`
	Views        int    `json:"views"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"-"`
}

func (a Article) TableName() string {
	return "article"
}

func (a Article) Count(db *gorm.DB) (int, error) {
	var count int
	if a.Title != "" {
		db.Where("title like ?", "%"+a.Title+"%")
	}
	if err := db.Model(&a).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (a Article) List(db *gorm.DB, pageOffset, pageSize int) ([]*Article, error) {
	var list []*Article
	if a.Title != "" {
		db.Where("title like ?", "%"+a.Title+"%")
	}
	err := db.Limit(pageSize).Offset(pageOffset).Find(&list).Error
	return list, err
}
