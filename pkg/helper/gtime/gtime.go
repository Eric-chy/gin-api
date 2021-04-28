package gtime

import (
	"bytes"
	"math"
	"strconv"
	"time"
)

// 修眠 秒
func SleepSecond(num int) {
	time.Sleep(time.Duration(num) * time.Second)
}

func GetNanoTime() int64 {
	times := time.Now().UnixNano()
	return times
}

func GetMicroTime() int64 {
	times := time.Now().UnixNano() / 1e3
	return times
}

/**
* @des 格式化时间转换函数
* @param timestamp int64 要转换的时间戳（秒）
* @return string
 */
func FormatTime(timestamp uint64) string {
	byTime := []uint64{365 * 24 * 60 * 60, 30 * 24 * 60 * 60, 7 * 24 * 60 * 60, 24 * 60 * 60, 60 * 60, 60, 1}
	unit := []string{"年前", "个月前", "星期前", "天前", "小时前", "分钟前", "秒前"}
	now := uint64(time.Now().Unix())
	ct := now - timestamp
	if ct <= 0 {
		return "刚刚"
	}
	var res string
	for i := 0; i < len(byTime); i++ {
		if ct < byTime[i] {
			continue
		}
		var temp = math.Floor(float64(ct / byTime[i]))
		ct = ct % byTime[i]
		if temp > 0 {
			var tempStr string
			tempStr = strconv.FormatFloat(temp, 'f', -1, 64)
			res = MergeString(tempStr, unit[i]) //此处调用了一个我自己封装的字符串拼接的函数（你也可以自己实现）
		}
		break //我想要的形式是精确到最大单位，即："2天前"这种形式，如果想要"2天12小时36分钟48秒前"这种形式，把此处break去掉，然后把字符串拼接调整下即可（别问我怎么调整，这如果都不会我也是无语）
	}
	return res
}

/**
* @des 拼接字符串
* @param args ...string 要被拼接的字符串序列
* @return string
 */
func MergeString(args ...string) string {
	buffer := bytes.Buffer{}
	for i := 0; i < len(args); i++ {
		buffer.WriteString(args[i])
	}
	return buffer.String()
}
