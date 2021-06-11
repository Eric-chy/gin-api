package dao

import (
	"gin-api/internal/model"
	"gin-api/pkg/app"
)

func (d *Dao) CountArticle(title string) (int, error) {
	var article model.Article
	article = model.Article{Title: title}
	return article.Count(d.engine)
}

func (d *Dao) GetArticleList(title string, page, pageSize int) ([]*model.Article, error) {
	var article model.Article
	article = model.Article{Title: title}
	pageOffset := app.GetPageOffset(page, pageSize)
	return article.List(d.engine, pageOffset, pageSize)
}

//func (d *Dao) CreateArticle(name string, state uint8, createdBy string) error {
//	tag := model.Article{
//		Model: &model.Model{CreatedAt: createdBy},
//	}
//
//	return tag.Create(d.engine)
//}

//func (d *Dao) UpdateArticle(id uint32, name string, state uint8, modifiedBy string) error {
//	tag := model.Tag{
//		Model: &model.Model{ID: id},
//	}
//	values := map[string]interface{}{
//		"state":       state,
//		"modified_by": modifiedBy,
//	}
//	if name != "" {
//		values["name"] = name
//	}
//
//	return tag.Update(d.engine, values)
//}

//func (d *Dao) DeleteArticle(id uint32) error {
//	tag := model.Tag{Model: &model.Model{ID: id}}
//	return tag.Delete(d.engine)
//}
