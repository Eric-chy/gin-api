package service

import (
	"ginpro/internal/model"
	"ginpro/pkg/app"
)

func (svc *Service) CountArticle(title string) (int, error) {
	return svc.dao.CountArticle(title)
}

func (svc *Service) GetArticleList(title string, pager *app.Pager) ([]*model.Article, error) {
	return svc.dao.GetArticleList(title, pager.Page, pager.PageSize)
}