package api

import (
	"fmt"
	"ginpro/common/dict"
	"ginpro/internal/service"
	"ginpro/pkg/app"
	"ginpro/pkg/gredis"
	"ginpro/pkg/helper/email"
	"ginpro/pkg/helper/gjson"
	"ginpro/pkg/httpclient"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type Article struct{}

func NewArticle() Article {

	return Article{}
}

// ArticleList
// @Summary 获取列表
// @Produce  json
// @Param name query string false "名称" maxlength(100)
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} model.Article"成功"
// @Failure 400 {object} dict.Error "请求错误"
// @Failure 500 {object} dict.Error "内部错误"
// @Router /api/articles [get]
func (a *Article) ArticleList(c *gin.Context) {
	param := struct {
		Title string `form:"title" binding:"max=100"`
	}{
		Title: c.Param("title"),
	}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		app.Error(c, dict.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	pager := app.Pager{Page: app.GetPage(c), PageSize: app.GetPageSize(c)}

	svc := service.New(c.Request.Context())
	totalRows, err := svc.CountArticle(param.Title)
	if err != nil {
		app.Error(c, dict.ErrGetArtCountFail)
		return
	}
	articles, err := svc.GetArticleList(param.Title, &pager)
	if err != nil {
		app.Error(c, dict.ErrGetArtListFail)
		return
	}
	for _, article := range articles {
		num, err := gredis.GetNum("article" + strconv.Itoa(int(article.Id)))
		if err != nil {
			num = 1
		}
		article.Views += num
	}
	app.SuccessList(c, articles, totalRows)
	return
}

func (a *Article) ArticleDetail(c *gin.Context) {
	json := make(map[string]interface{}) //注意该结构接受的内容
	c.BindJSON(&json)
	app.Success(c, map[string]interface{}{
		"name": json["name"],
		"age":  json["age"],
		"sex":  json["sex"],
	})

	//app.Error(c, dict.InvalidParams)
}

func (a *Article) SendEmail(c *gin.Context) {
	err := email.SendMail([]string{"test@126.com"}, "test", "测试一下邮件")
	if err != nil {
		fmt.Println(err)
		app.Error(c, dict.ServerError)
	}
	app.Success(c, map[string]string{"name": "test"})
	return
}

func (a *Article) Curl(c *gin.Context) {
	res, err := httpclient.New().Timeout(3*time.Second).Post("http://localhost:8088/api/articles/1", gin.H{"name": "aaa", "age": "12", "sex": "1"})

	if err != nil {
		fmt.Println(err)
	}
	defer res.Close()
	s := res.ReadAllString()
	r := gjson.JsonDecode(s)
	fmt.Println(r)
	fmt.Println(r["code"])
	fmt.Println(r["msg"])
	fmt.Println(r["data"])
	fmt.Println(r["elapsed"])
	fmt.Println(s)
	app.Success(c, r)
	return
}
