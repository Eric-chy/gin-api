package app

import (
	"gin-api/config"
	"gin-api/pkg/helper/convert"
	"github.com/gin-gonic/gin"
)

func GetPage(c *gin.Context) int {
	page := convert.Str(c.Query("page")).ToInt()
	if page <= 0 {
		return 1
	}

	return page
}

func GetPageSize(c *gin.Context) int {
	pageSize := convert.Str(c.Query("page_size")).ToInt()
	if pageSize <= 0 {
		return config.Conf.App.DefaultPageSize
	}
	if pageSize > config.Conf.App.MaxPageSize {
		return config.Conf.App.MaxPageSize
	}

	return pageSize
}

func GetPageOffset(page, pageSize int) int {
	result := 0
	if page > 0 {
		result = (page - 1) * pageSize
	}

	return result
}
