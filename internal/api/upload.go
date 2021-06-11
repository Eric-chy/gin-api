package api

import (
	"gin-api/common/dict"
	"gin-api/internal/service"
	"gin-api/pkg/app"
	"gin-api/pkg/helper/convert"
	"gin-api/pkg/helper/files"
	"github.com/gin-gonic/gin"
)

type Upload struct{}

func NewUpload() Upload {
	return Upload{}
}

func (u Upload) UploadFile(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		app.Error(c, dict.ServerError.WithDetails(err.Error()))
		return
	}

	fileType := convert.Str(c.PostForm("type")).ToInt()
	if fileHeader == nil || fileType <= 0 {
		app.Error(c, dict.InvalidParams)
		return
	}
	svc := service.New(c.Request.Context())
	fileInfo, err := svc.UploadFile(files.FileType(fileType), file, fileHeader)
	if err != nil {
		app.Error(c, dict.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}

	app.Success(c, gin.H{
		"file_access_url": fileInfo.AccessUrl,
	})
}
