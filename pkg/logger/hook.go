package logger

import (
	rotate "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"log"
)

/**
WithMaxAge和WithRotationCount二者只能设置一个
WithMaxAge设置文件清理前的最长保存时间
WithRotationCount设置文件清理前最多保存的个数
*/
func NewLfsHook(filePath string) logrus.Hook {
	infoWriter := newRotate("info", filePath)
	warnWriter := newRotate("warn", filePath)
	errorWriter := newRotate("error", filePath)
	writeMap := lfshook.WriterMap{
		logrus.DebugLevel: infoWriter,
		logrus.InfoLevel:  infoWriter,
		logrus.WarnLevel:  warnWriter,
		logrus.ErrorLevel: errorWriter,
		logrus.FatalLevel: errorWriter,
		logrus.PanicLevel: errorWriter,
	}
	lfsHook := lfshook.NewHook(writeMap, new(LogFormatter))
	return lfsHook
}

func newRotate(level, filePath string) *rotate.RotateLogs {
	writer, err := rotate.New(
		// 分割后的文件名称
		filePath + ".%Y%m%d." + level + ".log",
		// 生成软链，指向最新日志文件
		//rotate.WithLinkName(filePath),
		// 设置日志切割时间间隔(1天)
		//rotate.WithRotationTime(24*time.Hour),
		//rotate.WithMaxAge(30*24*time.Hour),
	)
	if err != nil {
		log.Println("rotate"+level+" err:", err)
	}
	return writer
}
