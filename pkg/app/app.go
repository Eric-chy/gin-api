package app

import (
	"bytes"
	"gin-api/common/dict"
	"gin-api/pkg/helper/gjson"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	Elapsed float64     `json:"elapsed"`
}

type Pager struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	TotalRows int `json:"total_rows"`
}

//Success 正常返回
func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = make([]string, 0)
	}
	response := Response{Code: 0, Msg: "success", Data: data, Elapsed: GetElapsed(c)}
	c.Set("responseData", response)
	c.JSON(http.StatusOK, response)
}

//SuccessList 分页返回
func SuccessList(c *gin.Context, list interface{}, totalRows int) {
	data := gin.H{
		"list": list,
		"pager": Pager{
			Page:      GetPage(c),
			PageSize:  GetPageSize(c),
			TotalRows: totalRows,
		},
	}
	e := dict.Success
	response := Response{Code: e.Code(), Msg: e.Msg(), Data: data, Elapsed: GetElapsed(c)}
	c.Set("responseData", response)
	c.JSON(http.StatusOK, response)
}

//Error 使用公共配置的消息返回
func Error(c *gin.Context, err *dict.Error) {
	response := Response{Code: err.Code(), Msg: err.Msg(), Elapsed: GetElapsed(c)}
	details := err.Details()
	if err.Level() == "" { //默认错误返回为warn，不记录日志到sentry
		err = err.WithLevel("warn")
	}
	SetLevel(c, err.Level())
	if len(details) > 0 {
		SetDetail(c, err.Details())
		if err.Level() != dict.LevelError {
			response.Data = details
		}
	}
	c.Set("responseData", response)
	c.JSON(err.StatusCode(), response)
}

func SetLevel(c *gin.Context, level interface{}) {
	c.Set("level", level)
}

func SetDetail(c *gin.Context, detail interface{}) {
	c.Set("detail", detail)
}

func GetLevel(c *gin.Context) interface{} {
	return Get(c, "level")
}

func GetDetail(c *gin.Context) interface{} {
	return Get(c, "detail")
}

func Get(c *gin.Context, key string) interface{} {
	val, _ := c.Get(key)
	return val
}

func GetElapsed(c *gin.Context) float64 {
	elapsed := 0.00
	if requestTime := Get(c, "beginTime"); requestTime != nil {
		elapsed = float64(time.Since(requestTime.(time.Time))) / 1e9
	}
	return elapsed
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
