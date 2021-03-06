# Gin项目

#### 介绍
<h6 style="color:red;font-weight:100;">有兴趣请star一下，以下是基于gin开发的项目接口，将持续更新，本项目包含mysql，redis，elasticsearch，mongo，rabbitmq，kafka，jaeger，单机限流，分布式限流，sentry, jwt，请求参数验证，发送邮件，图片上传，httpclient用于请求第三方接口等, cmd目录下执行```go run genModel.go```可自动生成model文件，后面会补上grpc的部分，另外可以关注我的博客http://www.cyz.show 以下3-6所有组件的安装可以参考我的博客：http://www.cyz.show/archives/143/</h6>

#### 目录结构
~~~
ginpro  根目录
├─boot  初始化启动数据库连接等
├─cmd  自动生成model文件
├─common  通用数据字典和全局变量
│   ├─dict   数据字典，错误码和常用参数
│   └─global 全局变量    
├─config 系统配置文件目录
│   ├─config.go  配置初始化
│   ├─dev.yaml   开发机配置
│   └─qa.yaml   测试环境配置    
├─docs  swagger文档目录（下面三个文件在根目录swag init命令生成）
│  ├─docs.go            
│  ├─swagger.json            
│  └─swagger.yaml
├─internal  
│  ├─api  接口                    
│  ├─dao dao层，对数据库的增删改查            
│  ├─middleware 中间件            
│  ├─model model层，数据库字段表名等            
│  ├─router 路由            
│  └─service 
├─pkg  
│  ├─app  接口     
│  │  ├─app.go 接口响应等方法封装        
│  │  ├─form.go 表单验证封装            
│  │  ├─jwt.go jwt鉴权            
│  │  └─pagination.go 分页  
│  ├─es elasticsearch     
│  ├─gredis  redis    
│  ├─helper
│  │  ├─convert 常用转换        
│  │  ├─email 邮件发送            
│  │  ├─files 文件操作相关            
│  │  ├─gjson json操作            
│  │  └─gtime 时间相关操作  
│  ├─httpclient 请求第三方，类似于curl   
│  ├─limiter 限流            
│  ├─logger 日志                   
│  ├─mgodb mongodb                   
│  ├─rabbitmq rabbitmq                   
│  ├─security md5，aes加密等                       
│  └─tracer 链路追踪
├─storage               
│  ├─logs 日志            
│  └─uploads  上传的文件 
├─go.mod   模块管理   
└─main.go 入口文件
~~~

#### 启动教程
1. 安装mysql
2. 安装jaeger:
   ```dockerfile
   docker run -d --name jaeger \
   -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
   -p 5775:5775/udp \
   -p 6831:6831/udp \
   -p 6832:6832/udp \
   -p 5778:5778 \
   -p 16686:16686 \
   -p 14268:14268 \
   -p 9411:9411 \
   jaegertracing/all-in-one:latest
   ```
3.  安装es（以下3-6都可以通过在boot/boot.go的InitApp方法中选择决定初始化是否安装，如果不需要可以注释掉boot/boot.go中对应启动代码或者在配置文件中加上开关配置参数判断是否需要启动）
4.  安装redis
5.  安装mongo
6.  安装rabbitmq
7.  安装sentry,如果不想安装，请修改pkg/logger/logrus.go文件中的下面代码注释掉,或者自己在配置文件中设置开关，安装可以参考我的博客：http://www.cyz.show/archives/252/
    ```golang
    hook, err := logrus_sentry.NewSentryHook(config.Conf.Sentry.Dsn, []logrus.Level{
    logrus.PanicLevel,
    logrus.FatalLevel,
    logrus.ErrorLevel,
    })
    if err == nil {
    global.Logger.Hooks.Add(hook)
    hook.Timeout = 0
    hook.StacktraceConfiguration.Enable = true
    }
    ```
8.  在项目根目录下执行```go mod tidy```，不清楚的请自行百度一下，先设置好proxy再执行命令，如：GOPROXY=https://goproxy.cn
9. 为了减少手动写model文件的麻烦，cmd目录提供了自动生成model文件的工具，使用如下：
   在cmd目录下执行```go run genModel.go``` 加上以下参数或者不加
   ```
   -c string 指定要使用的配置文件路径 (default "../config/")
   -r string 是否替换旧文件生成 (default "n" "n|y")
   -d string 数据库名,不填则按配置文件来
   -f string 指定要使用的配置文件名 (default "dev") 使用dev.yaml
   -m string 指定要生成的model路径 (default "../internal/model/")
   -t string 表名，多个使用英文半角,分割，不填则生成数据库下所有表的model
   ```
10.  根目录```go run main.go```即可启动
11.  为了方便开发，一般使用热更新，安装fresh，在根目录下执行```go get github.com/pilu/fresh```，然后使用fresh命令即可启动，和上面第10步骤二选一
12. swagger安装，生成文档，非必需
    ```
    go get -u github.com/swaggo/swag/cmd/swag@v1.6.5 
    go get -u github.com/swaggo/gin-swagger@v1.2.0
    go get -u github.com/swaggo/files
    go get -u github.com/alecthomas/template
    ```
    验证是否安装成功： swag -v
    对controller即接口进行注解
    ```
    // @Summary 获取列表
    // @Produce  json
    // @Param name query string false "名称" maxlength(100)
    // @Param state query int false "状态" Enums(0, 1) default(1)
    // @Param page query int false "页码"
    // @Param page_size query int false "每页数量"
    // @Success 200 {object} model.Tag "成功"
    // @Failure 400 {object} code.Error "请求错误"
    // @Failure 500 {object} code.Error "内部错误"
    // @Router /api/list [get]
    func (c *Controller) List (c *gin.Context) {
        app.Success(c, nil)
    }
    区分项目的话，在main入口函数添加注解：
    // @title gin系统
    // @version 1.0
    // @description gin开发的系统
    // @termsOfService 
    func main(){}
    在model文件中添加
    type ArticleSwagger struct {
        List  []*Article
        Pager *app.Pager
    }
    ```
    设置路由，在apiRouter.go中先(否则会报:Failed to load spec)
    ```
    import(
        _ "ginpro/docs"
    )
    ```
    再设置
    ```
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    ```
    生成swagger文档：在根目录执行swag init
    swagger文档查看 http://127.0.0.1:8001/swagger/index.html查看
    
13.如果不需要其中的某些组件，如es，redis，mongo等，可以在boot/boot.go init方法中注释掉相关的即可，或者在配置文件中设置开关（自行实现即可）
#### 使用说明

1.  启动项目后我们可以看到对应的路由信息，可以使用浏览器或者postman之类的进行访问
2.  链路信息查看：http://127.0.0.1:16686/
3.  文件上传测试 ```curl -X POST http://127.0.0.1:8000/upload/file -F file=@{file_path} -F type=1```
4. 项目启动后访问 http://127.0.0.1:8088/api/articles
5.  建表语句，目前只是简单展示，所以只建立了一个简单的表
```mysql
CREATE DATABASE blog;
USE blog;
CREATE TABLE `article` (
`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
`title` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
`introduction` varchar(500) COLLATE utf8mb4_unicode_ci,
`views` int(11) NOT NULL DEFAULT '0',
`content` varchar(5000) COLLATE utf8mb4_unicode_ci,
`created_at` timestamp NULL DEFAULT NULL,
`updated_at` timestamp NULL DEFAULT NULL,
PRIMARY KEY (`id`)
) ENGINE=Innodb DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
INSERT INTO blog.article VALUES(NULL, "我的第一篇文章", "文章简介", 100, "文章的内容很好看", "2020-02-02 02:22:22", "2020-02-02 02:22:22") 
```
