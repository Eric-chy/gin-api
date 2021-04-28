package logger

import (
	"fmt"
	"ginpro/common/global"
	"ginpro/pkg/helper/convert"
	"ginpro/pkg/helper/gjson"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"time"
)

//日志自定义格式
type LogFormatter struct{}

//格式详情
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
		"%s %s [%s] [GID:%d] [RID:%d] #file:%s:%d #msg:%s #content:%v\n",
		timestamp, global.Ip, level, convert.GID(), global.RequestID, file, line, entry.Message, content,
	)
	return []byte(msg), nil
}
