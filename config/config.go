package config

import (
	"flag"
	"github.com/spf13/viper"
	"time"
)

var Conf Yaml

type Yaml struct {
	App struct {
		RunMode              string        `yaml:"runMode"`
		Port                 string        `yaml:"port"`
		ReadTimeout          time.Duration `yaml:"readTimeout"`
		WriteTimeout         time.Duration `yaml:"writeTimeout"`
		AppName              string        `yaml:"appName"`
		LogDir               string        `yaml:"logDir"`
		AESKey               string        `yaml:"aesKey"`
		MaxPageSize          int           `yaml:"maxPageSize"`
		DefaultPageSize      int           `yaml:"defaultPageSize"`
		UploadImageMaxSize   int           `yaml:"defaultPageSize"`
		UploadSavePath       string        `yaml:"uploadSavePath"`
		UploadServerUrl      string        `yaml:"uploadServerUrl"`
		UploadImageAllowExts []string      `yaml:"uploadImageAllowExts"`
	}
	Database struct {
		Driver      string        `yaml:"driver"`
		Protocol    string        `yaml:"tcp"`
		Host        string        `yaml:"host"`
		Port        int           `yaml:"port"`
		User        string        `yaml:"user"`
		Password    string        `yaml:"password"`
		Name        string        `yaml:"name"`
		Prefix      string        `yaml:"prefix"`
		RunMode     string        `yaml:"runMode"`
		MaxIdles    int           `yaml:"maxIdles"`
		MaxOpens    int           `yaml:"maxOpens"`
		MaxLifetime time.Duration `yaml:"maxLifetime"`
	}
	Redis struct {
		Driver      string        `yaml:"driver"`
		Protocol    string        `yaml:"tcp"`
		Host        string        `yaml:"host"`
		Port        string        `yaml:"port"`
		Password    string        `yaml:"password"`
		MaxIdle     int           `yaml:"maxIdle"`
		MaxActive   int           `yaml:"maxActive"`
		IdleTimeout time.Duration `yaml:"IdleTimeout"`
	}
	Sentry struct {
		Dsn string `yaml:"dsn"`
	}
	Jwt struct {
		Secret string        `yaml:"secret"`
		Issuer string        `yaml:"issuer"`
		Expire time.Duration `yaml:"expire"`
	}
	Jaeger struct {
		Link string `yaml:"link"`
	}
	Email struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		UserName string `yaml:"userName"`
		Password string `yaml:"password"`
		IsSSL    bool   `yaml:"isSSL"`
		From     string `yaml:"from"`
	}
	Es struct {
		Link                string `yaml:"link"`
		MaxIdleConnsPerHost int    `yaml:"maxIdleConnsPerHost"`
	}
	Mongo struct {
		Host        string        `yaml:"host"`
		MaxPoolSize uint64        `yaml:"maxPoolSize"`
		Timeout     time.Duration `yaml:"timeout"`
	}
}

var (
	port    string
	runMode string
	cfg     string
	env     string
)
var vp *viper.Viper

func Init() {
	setupFlag()
	vp = viper.New()
	vp.AddConfigPath(cfg)
	vp.SetConfigName("dev")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		panic(err)
	}
	//直接整个解析
	err = vp.Unmarshal(&Conf)
	if err != nil {
		panic(err)
	}
}

func setupFlag() {
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&cfg, "config", "./config/", "指定要使用的配置文件路径")
	flag.StringVar(&env, "env", "dev", "指定要使用的配置文件")
	flag.Parse()
}
