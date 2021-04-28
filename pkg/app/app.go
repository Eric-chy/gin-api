package app

import (
	"bytes"
	"ginpro/common/dict"
	"ginpro/pkg/helper/gjson"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Pager struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	TotalRows int `json:"total_rows"`
}

//Success 正常返回
func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, data)
}

//SuccessList 分页返回
func SuccessList(c *gin.Context, list interface{}, totalRows int) {
	c.JSON(http.StatusOK, gin.H{
		"list": list,
		"pager": Pager{
			Page:      GetPage(c),
			PageSize:  GetPageSize(c),
			TotalRows: totalRows,
		},
	})
}

//Error 使用公共配置的消息返回
func Error(c *gin.Context, err *dict.Error) {
	response := gin.H{"code": err.Code(), "msg": err.Msg()}
	details := err.Details()
	if err.Level() == "" {//默认错误返回为warn，不记录日志到sentry
		err = err.WithLevel("warn")
	}
	SetLevel(c, err.Level())
	if len(details) > 0 {
		SetDetail(c, err.Details())
		if err.Level() != dict.LevelError {
			response["detail"] = details
		}
	}
	c.JSON(err.StatusCode(), response)
}

func SetResponseData(c *gin.Context, response interface{}) {
	c.Set("responseData", response)
}

func SetLevel(c *gin.Context, level interface{}) {
	c.Set("level", level)
}

func GetLevel(c *gin.Context) interface{} {
	level, _ := c.Get("level")
	return level
}

func SetDetail(c *gin.Context, detail interface{}) {
	c.Set("detail", detail)
}

func GetDetail(c *gin.Context) interface{} {
	detail, _ := c.Get("detail")
	return detail
}

func JsonParams(c *gin.Context) map[string]interface{} {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}
	// 将取出来的body内容重新插入body，否则ShouldBindJSON无法绑定参数
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	return gjson.JsonDecode(string(b))
}