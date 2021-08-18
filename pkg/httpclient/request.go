//Package httpclient 一个简单的http客户端，用于请求第三方接口，简单实现了GET和POST方法，比较粗糙，一般情况下够用
package httpclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//Post domain string "https://aaa.com"
//Post data url.values, eg:
//data := url.Values{}
//data.Add("Name", "Ewan")
//data.Add("Age", 12)
func Post(domain string, data url.Values) (*Response, error) {
	//增加header选项
	r, _ := http.NewRequest("POST", domain, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return response(resp)
}

func Get(domain string, data url.Values) (*Response, error) {
	uri, _ := url.Parse(domain)
	uri.RawQuery = data.Encode()
	//r, _ := http.NewRequest("GET", domain + data.Encode(), nil)
	r, _ := http.NewRequest("GET", uri.String(), nil)
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return response(resp)
}

func response(resp *http.Response) (*Response, error) {
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(fmt.Sprintf("http error,http code: %d,body:%s", resp.StatusCode, body))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
