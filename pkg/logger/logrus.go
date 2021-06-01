package logger

import (
	"ginpro/common/global"
	"ginpro/config"
	"ginpro/pkg/helper/files"
	"github.com/sirupsen/logrus"
)

func Init() {
	global.Logger = logrus.New()
	if config.Conf.Sentry.Dsn != "" {
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
	}
	// 设置日志格式为json格式
	global.Logger.SetFormatter(&logrus.JSONFormatter{})
	//设置文件输出
	f, logFilePath := files.LogFile()
	// 日志消息输出可以是任意的io.writer类型，这里我们获取文件句柄，将日志输出到文件
	global.Logger.SetOutput(f)
	// 设置日志级别为debug以上
	global.Logger.SetLevel(logrus.DebugLevel)
	// 设置显示文件名和行号
	global.Logger.SetReportCaller(true)
	// 设置rotatelogs日志分割Hook
	global.Logger.AddHook(NewLfsHook(logFilePath))

}
