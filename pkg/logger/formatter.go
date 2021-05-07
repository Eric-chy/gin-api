package logger

import (
	"fmt"
	"ginpro/pkg/helper/convert"
	"ginpro/pkg/helper/gjson"
	"ginpro/pkg/helper/gtime"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"time"
)

//LogFormatter 日志自定义格式
type LogFormatter struct{}

//Format 格式详情
func (s *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	var file string
	var line int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		line = entry.Caller.Line
	}

	level := strings.ToUpper(entry.Level.String())
	content := gjson.JsonEncode(entry.Data)
	msg := fmt.Sprintf(
		"【%s】 [%s] [ip:%s] [GID:%d] [RID:%d]\r\n#File:%s:%d \r\n#Msg:%s \r\n#Content:%v\n",
		level, timestamp, entry.Data["ip"], convert.GID(), gtime.GetMicroTime(), file, line, entry.Message, content,
	)
	return []byte(msg), nil
}
