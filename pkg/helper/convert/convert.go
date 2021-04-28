package convert

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"net"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"
)
type Str string
//StructToMap 结构体转 map[string]interface{}
func StructToMap(in interface{}, tagName string) map[string]interface{} {
	t := reflect.TypeOf(in)
	v := reflect.ValueOf(in)

	maps := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		// 指定tagName值为map中key
		k := f.Tag.Get(tagName)
		v := v.FieldByName(f.Name).Interface()
		maps[k] = v
	}
	return maps
}

func TypeAssertion(i interface{}) string {
	var key string
	if i == nil {
		return key
	}

	switch i.(type) {
	case float64:
		ft := i.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := i.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := i.(int)
		key = strconv.Itoa(it)
	case uint:
		it := i.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := i.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := i.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := i.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := i.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := i.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := i.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := i.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := i.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = i.(string)
	case []byte:
		key = string(i.([]byte))
	default:
		newValue, _ := json.Marshal(i)
		key = string(newValue)
	}

	return key
}

//RandomString 生成随机字符串
func RandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(letters))]
	}
	return string(b)
}

//GID 获取当前协程id
func GID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// Ip2Long 把ip字符串转为数值
func Ip2Long(ip string) (uint64, error) {
	b := net.ParseIP(ip).To4()
	if b == nil {
		return 0, errors.New("invalid ipv4 format")
	}
	return uint64(b[3]) | uint64(b[2])<<8 | uint64(b[1])<<16 | uint64(b[0])<<24, nil
}

// Long2Ip 把数值转为ip字符串
func Long2Ip(i uint64) (string, error) {
	if i > math.MaxUint32 {
		return "", errors.New("beyond the scope of ipv4")
	}

	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)

	return ip.String(), nil
}

func (s Str) String() string {
	return string(s)
}

func (s Str) Int() (int, error) {
	v, err := strconv.Atoi(s.String())
	return v, err
}

func (s Str) ToInt() int {
	v, _ := s.Int()
	return v
}

func (s Str) UInt32() (uint32, error) {
	v, err := strconv.Atoi(s.String())
	return uint32(v), err
}

func (s Str) ToUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func UnderscoreToUpperCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	return strings.Replace(s, " ", "", -1)
}

func UnderscoreToLowerCamelCase(s string) string {
	s = UnderscoreToUpperCamelCase(s)
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

func CamelCaseToUnderscore(s string) string {
	var output []rune
	for i, r := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(r))
			continue
		}
		if unicode.IsUpper(r) {
			output = append(output, '_')
		}
		output = append(output, unicode.ToLower(r))
	}
	return string(output)
}