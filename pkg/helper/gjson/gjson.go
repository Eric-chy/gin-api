package gjson

import (
	"encoding/json"
)

func JsonEncode(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func JsonDecode(s string) map[string]interface{} {
	data := make(map[string]interface{})
	if s == "" {
		return data
	}
	if err := json.Unmarshal([]byte(s), &data); err != nil {
		panic(err)
	}
	return data
}
