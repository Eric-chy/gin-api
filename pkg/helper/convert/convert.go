package convert

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unsafe"
)

type Str string

var (
	// DefaultTrimChars are the characters which are stripped by Trim* functions in default.
	DefaultTrimChars = string([]byte{
		'\t', // Tab.
		'\v', // Vertical tab.
		'\n', // New line (line feed).
		'\r', // Carriage return.
		'\f', // New page.
		' ',  // Ordinary space.
		0x00, // NUL-byte.
		0x85, // Delete.
		0xA0, // Non-breaking space.
	})
)

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

type apiString interface {
	String() string
}

type apiError interface {
	Error() string
}

func String(any interface{}) string {
	if any == nil {
		return ""
	}
	switch value := any.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	case time.Time:
		if value.IsZero() {
			return ""
		}
		return value.String()
	case *time.Time:
		if value == nil {
			return ""
		}
		return value.String()
	default:
		// Empty checks.
		if value == nil {
			return ""
		}
		if f, ok := value.(apiString); ok {
			// If the variable implements the String() interface,
			// then use that interface to perform the conversion
			return f.String()
		}
		if f, ok := value.(apiError); ok {
			// If the variable implements the Error() interface,
			// then use that interface to perform the conversion
			return f.Error()
		}
		// Reflect checks.
		var (
			rv   = reflect.ValueOf(value)
			kind = rv.Kind()
		)
		switch kind {
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return ""
			}
		case reflect.String:
			return rv.String()
		}
		if kind == reflect.Ptr {
			return String(rv.Elem().Interface())
		}
		// Finally we use json.Marshal to convert.
		if jsonContent, err := json.Marshal(value); err != nil {
			return fmt.Sprint(value)
		} else {
			return string(jsonContent)
		}
	}
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

func Trim(str string, characterMask ...string) string {
	trimChars := DefaultTrimChars
	if len(characterMask) > 0 {
		trimChars += characterMask[0]
	}
	return strings.Trim(str, trimChars)
}

// Contains reports whether <substr> is within <str>, case-sensitively.
func Contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

func SplitAndTrim(str, delimiter string, characterMask ...string) []string {
	array := make([]string, 0)
	for _, v := range strings.Split(str, delimiter) {
		v = Trim(v, characterMask...)
		if v != "" {
			array = append(array, v)
		}
	}
	return array
}

func Map(value interface{}) map[string]interface{} {
	if value == nil {
		return nil
	}
	// Assert the common combination of types, and finally it uses reflection.
	dataMap := make(map[string]interface{})
	switch r := value.(type) {
	case string:
		// If it is a JSON string, automatically unmarshal it!
		if len(r) > 0 && r[0] == '{' && r[len(r)-1] == '}' {
			if err := json.Unmarshal([]byte(r), &dataMap); err != nil {
				return nil
			}
		} else {
			return nil
		}
	case []byte:
		// If it is a JSON string, automatically unmarshal it!
		if len(r) > 0 && r[0] == '{' && r[len(r)-1] == '}' {
			if err := json.Unmarshal(r, &dataMap); err != nil {
				return nil
			}
		} else {
			return nil
		}
	case map[interface{}]interface{}:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[interface{}]string:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[interface{}]int:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[interface{}]uint:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[interface{}]float32:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[interface{}]float64:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[string]bool:
		for k, v := range r {
			dataMap[k] = v
		}
	case map[string]int:
		for k, v := range r {
			dataMap[k] = v
		}
	case map[string]uint:
		for k, v := range r {
			dataMap[k] = v
		}
	case map[string]float32:
		for k, v := range r {
			dataMap[k] = v
		}
	case map[string]float64:
		for k, v := range r {
			dataMap[k] = v
		}
	case map[string]interface{}:
		return r
	case map[int]interface{}:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[int]string:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[uint]string:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	}
	return dataMap
}

func UnsafeBytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
