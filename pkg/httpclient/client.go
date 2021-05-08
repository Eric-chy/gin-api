//Package httpclient  移植自gf
package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"ginpro/config"
	"ginpro/pkg/helper/convert"
	"ginpro/pkg/helper/files"
	"ginpro/pkg/helper/gregex"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	clientMiddlewareKey = "__clientMiddlewareKey"
	fileUploadingKey    = "@file:"
)

// Response is the struct for client request response.
type Response struct {
	*http.Response
	request     *http.Request
	requestBody []byte
	cookies     map[string]string
}

// HandlerFunc middleware handler func
type HandlerFunc = func(c *Client, r *http.Request) (*Response, error)

// Client is the HTTP client for HTTP request management.
type Client struct {
	http.Client                         // Underlying HTTP Client.
	ctx               context.Context   // Context for each request.
	dump              bool              // Mark this request will be dumped.
	parent            *Client           // Parent http client, this is used for chaining operations.
	header            map[string]string // Custom header map.
	cookies           map[string]string // Custom cookie map.
	prefix            string            // Prefix for request.
	authUser          string            // HTTP basic authentication: user.
	authPass          string            // HTTP basic authentication: pass.
	retryCount        int               // Retry count when request fails.
	retryInterval     time.Duration     // Retry interval when request fails.
	middlewareHandler []HandlerFunc     // Interceptor handlers
}

var (
	defaultClientAgent = fmt.Sprintf(`GinHTTPClient %s`, config.Conf.App.Version)
)

// New creates and returns a new HTTP client object.
func New() *Client {
	c := &Client{
		Client: http.Client{
			Transport: &http.Transport{
				// No validation for https certification of the server in default.
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				DisableKeepAlives: true,
			},
		},
		header:  make(map[string]string),
		cookies: make(map[string]string),
	}
	c.header["User-Agent"] = defaultClientAgent
	return c
}

func (c *Client) Timeout(t time.Duration) *Client {
	newClient := c
	if c.parent == nil {
		newClient = c.Clone()
	}
	newClient.SetTimeout(t)
	return newClient
}

// Clone deeply clones current client and returns a new one.
func (c *Client) Clone() *Client {
	newClient := New()
	*newClient = *c
	newClient.header = make(map[string]string)
	newClient.cookies = make(map[string]string)
	for k, v := range c.header {
		newClient.header[k] = v
	}
	for k, v := range c.cookies {
		newClient.cookies[k] = v
	}
	return newClient
}

//SetTimeout sets the request timeout for the client.
func (c *Client) SetTimeout(t time.Duration) *Client {
	c.Client.Timeout = t
	return c
}

// Get send GET request and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Get(url string, data ...interface{}) (*Response, error) {
	return c.DoRequest("GET", url, data...)
}

func (c *Client) GetContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("GET", url, data...))
}

func (c *Client) RequestBytes(method string, url string, data ...interface{}) []byte {
	resp, err := c.DoRequest(method, url, data...)
	if err != nil {
		return nil
	}
	defer resp.Close()
	return resp.ReadAll()
}

func (r *Response) Close() error {
	if r == nil || r.Response == nil || r.Response.Close {
		return nil
	}
	r.Response.Close = true
	return r.Response.Body.Close()
}

// Put send PUT request and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Put(url string, data ...interface{}) (*Response, error) {
	return c.DoRequest("PUT", url, data...)
}

// Post sends request using HTTP method POST and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Post(url string, data ...interface{}) (*Response, error) {
	return c.DoRequest("POST", url, data...)
}

func (c *Client) PostContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("POST", url, data...))
}

// Delete send DELETE request and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Delete(url string, data ...interface{}) (*Response, error) {
	return c.DoRequest("DELETE", url, data...)
}

// Head send HEAD request and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Head(url string, data ...interface{}) (*Response, error) {
	return c.DoRequest("HEAD", url, data...)
}

// Patch send PATCH request and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Patch(url string, data ...interface{}) (*Response, error) {
	return c.DoRequest("PATCH", url, data...)
}

// Connect send CONNECT request and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Connect(url string, data ...interface{}) (*Response, error) {
	return c.DoRequest("CONNECT", url, data...)
}

// Options send OPTIONS request and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Options(url string, data ...interface{}) (*Response, error) {
	return c.DoRequest("OPTIONS", url, data...)
}

// Trace send TRACE request and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Trace(url string, data ...interface{}) (*Response, error) {
	return c.DoRequest("TRACE", url, data...)
}

// Use adds one or more middleware handlers to client.
func (c *Client) Use(handlers ...HandlerFunc) *Client {
	c.middlewareHandler = append(c.middlewareHandler, handlers...)
	return c
}

// Next calls next middleware.
// This is should only be call in HandlerFunc.
func (c *Client) Next(req *http.Request) (*Response, error) {
	if v := req.Context().Value(clientMiddlewareKey); v != nil {
		if m, ok := v.(*clientMiddleware); ok {
			return m.Next(req)
		}
	}
	return c.callRequest(req)
}

func (c *Client) Retry(retryCount int, retryInterval time.Duration) *Client {
	newClient := c
	if c.parent == nil {
		newClient = c.Clone()
	}
	newClient.SetRetry(retryCount, retryInterval)
	return newClient
}

func (c *Client) SetRetry(retryCount int, retryInterval time.Duration) *Client {
	c.retryCount = retryCount
	c.retryInterval = retryInterval
	return c
}

func (c *Client) PostBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("POST", url, data...)
}

func (c *Client) GetBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("GET", url, data...)
}

func (c *Client) DoRequest(method, url string, data ...interface{}) (resp *Response, err error) {
	req, err := c.prepareRequest(method, url, data...)
	if err != nil {
		return nil, err
	}

	// Client middleware.
	if len(c.middlewareHandler) > 0 {
		mdlHandlers := make([]HandlerFunc, 0, len(c.middlewareHandler)+1)
		mdlHandlers = append(mdlHandlers, c.middlewareHandler...)
		mdlHandlers = append(mdlHandlers, func(cli *Client, r *http.Request) (*Response, error) {
			return cli.callRequest(r)
		})
		ctx := context.WithValue(req.Context(), clientMiddlewareKey, &clientMiddleware{
			client:       c,
			handlers:     mdlHandlers,
			handlerIndex: -1,
		})
		req = req.WithContext(ctx)
		resp, err = c.Next(req)
	} else {
		resp, err = c.callRequest(req)
	}
	return resp, err
}

// prepareRequest verifies request parameters, builds and returns http request.
func (c *Client) prepareRequest(method, url string, data ...interface{}) (req *http.Request, err error) {
	method = strings.ToUpper(method)
	if len(c.prefix) > 0 {
		url = c.prefix + convert.Trim(url)
	}
	var params string
	if len(data) > 0 {
		switch c.header["Content-Type"] {
		case "application/json":
			switch data[0].(type) {
			case string, []byte:
				params = convert.String(data[0])
			default:
				if b, err := json.Marshal(data[0]); err != nil {
					return nil, err
				} else {
					params = string(b)
				}
			}
		case "application/xml":
			switch data[0].(type) {
			case string, []byte:
				params = convert.String(data[0])
			default:
				return nil, errors.New("xml error")
			}
		default:
			params = BuildParams(data[0])
		}
	}
	if method == "GET" {
		var bodyBuffer *bytes.Buffer
		if params != "" {
			switch c.header["Content-Type"] {
			case
				"application/json",
				"application/xml":
				bodyBuffer = bytes.NewBuffer([]byte(params))
			default:
				// It appends the parameters to the url
				// if http method is GET and Content-Type is not specified.
				if convert.Contains(url, "?") {
					url = url + "&" + params
				} else {
					url = url + "?" + params
				}
				bodyBuffer = bytes.NewBuffer(nil)
			}
		} else {
			bodyBuffer = bytes.NewBuffer(nil)
		}
		if req, err = http.NewRequest(method, url, bodyBuffer); err != nil {
			return nil, err
		}
	} else {
		if strings.Contains(params, "@file:") {
			// File uploading request.
			var (
				buffer = bytes.NewBuffer(nil)
				writer = multipart.NewWriter(buffer)
			)
			for _, item := range strings.Split(params, "&") {
				array := strings.Split(item, "=")
				if len(array[1]) > 6 && strings.Compare(array[1][0:6], "@file:") == 0 {
					path := array[1][6:]
					if files.CheckSavePath(path) {
						return nil, errors.New(fmt.Sprintf(`"%s" does not exist`, path))
					}
					if file, err := writer.CreateFormFile(array[0], files.Basename(path)); err == nil {
						if f, err := os.Open(path); err == nil {
							if _, err = io.Copy(file, f); err != nil {
								if err := f.Close(); err != nil {
									log.Errorf(`%+v`, err)
								}
								return nil, err
							}
							if err := f.Close(); err != nil {
								log.Errorf(`%+v`, err)
							}
						} else {
							return nil, err
						}
					} else {
						return nil, err
					}
				} else {
					if err = writer.WriteField(array[0], array[1]); err != nil {
						return nil, err
					}
				}
			}
			// Close finishes the multipart message and writes the trailing
			// boundary end line to the output.
			if err = writer.Close(); err != nil {
				return nil, err
			}

			if req, err = http.NewRequest(method, url, buffer); err != nil {
				return nil, err
			} else {
				req.Header.Set("Content-Type", writer.FormDataContentType())
			}
		} else {
			// Normal request.
			paramBytes := []byte(params)
			if req, err = http.NewRequest(method, url, bytes.NewReader(paramBytes)); err != nil {
				return nil, err
			} else {
				if v, ok := c.header["Content-Type"]; ok {
					// Custom Content-Type.
					req.Header.Set("Content-Type", v)
				} else if len(paramBytes) > 0 {
					if (paramBytes[0] == '[' || paramBytes[0] == '{') && json.Valid(paramBytes) {
						// Auto detecting and setting the post content format: JSON.
						req.Header.Set("Content-Type", "application/json")
					} else if gregex.IsMatchString(`^[\w\[\]]+=.+`, params) {
						// If the parameters passed like "name=value", it then uses form type.
						req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					}
				}
			}
		}
	}

	// Context.
	if c.ctx != nil {
		req = req.WithContext(c.ctx)
	} else {
		req = req.WithContext(context.Background())
	}
	// Custom header.
	if len(c.header) > 0 {
		for k, v := range c.header {
			req.Header.Set(k, v)
		}
	}
	// It's necessary set the req.Host if you want to custom the host value of the request.
	// It uses the "Host" value from header if it's not empty.
	if host := req.Header.Get("Host"); host != "" {
		req.Host = host
	}
	// Custom Cookie.
	if len(c.cookies) > 0 {
		headerCookie := ""
		for k, v := range c.cookies {
			if len(headerCookie) > 0 {
				headerCookie += ";"
			}
			headerCookie += k + "=" + v
		}
		if len(headerCookie) > 0 {
			req.Header.Set("Cookie", headerCookie)
		}
	}
	// HTTP basic authentication.
	if len(c.authUser) > 0 {
		req.SetBasicAuth(c.authUser, c.authPass)
	}
	return req, nil
}

func BuildParams(params interface{}, noUrlEncode ...bool) (encodedParamStr string) {
	// If given string/[]byte, converts and returns it directly as string.
	switch v := params.(type) {
	case string, []byte:
		return convert.String(params)
	case []interface{}:
		if len(v) > 0 {
			params = v[0]
		} else {
			params = nil
		}
	}
	// Else converts it to map and does the url encoding.
	m, urlEncode := convert.Map(params), true
	if len(m) == 0 {
		return convert.String(params)
	}
	if len(noUrlEncode) == 1 {
		urlEncode = !noUrlEncode[0]
	}
	// If there's file uploading, it ignores the url encoding.
	if urlEncode {
		for k, v := range m {
			if convert.Contains(k, fileUploadingKey) || convert.Contains(convert.String(v), fileUploadingKey) {
				urlEncode = false
				break
			}
		}
	}
	s := ""
	for k, v := range m {
		if len(encodedParamStr) > 0 {
			encodedParamStr += "&"
		}
		s = convert.String(v)
		if urlEncode && len(s) > 6 && strings.Compare(s[0:6], fileUploadingKey) != 0 {
			s = url.QueryEscape(s)
		}
		encodedParamStr += k + "=" + s
	}
	return
}

// callRequest sends request with give http.Request, and returns the responses object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) callRequest(req *http.Request) (resp *Response, err error) {
	resp = &Response{
		request: req,
	}
	// Dump feature.
	// The request body can be reused for dumping
	// raw HTTP request-response procedure.
	if c.dump {
		reqBodyContent, _ := ioutil.ReadAll(req.Body)
		resp.requestBody = reqBodyContent
		req.Body = NewReadCloser(reqBodyContent, false)
	}
	for {
		if resp.Response, err = c.Do(req); err != nil {
			// The response might not be nil when err != nil.
			if resp.Response != nil {
				if err := resp.Response.Body.Close(); err != nil {
					log.Errorf(`%+v`, err)
				}
			}
			if c.retryCount > 0 {
				c.retryCount--
				time.Sleep(c.retryInterval)
			} else {
				//return resp, err
				break
			}
		} else {
			break
		}
	}
	return resp, err
}

func (c *Client) ContentXml() *Client {
	newClient := c
	if c.parent == nil {
		newClient = c.Clone()
	}
	newClient.SetContentType("application/xml")
	return newClient
}

// SetContentType sets HTTP content type for the client.
func (c *Client) SetContentType(contentType string) *Client {
	c.header["Content-Type"] = contentType
	return c
}

func (c *Client) ContentJson() *Client {
	newClient := c
	if c.parent == nil {
		newClient = c.Clone()
	}
	newClient.SetContentType("application/json")
	return newClient
}

func (c *Client) ContentType(contentType string) *Client {
	newClient := c
	if c.parent == nil {
		newClient = c.Clone()
	}
	newClient.SetContentType(contentType)
	return newClient
}

func (c *Client) Header(m map[string]string) *Client {
	newClient := c
	if c.parent == nil {
		newClient = c.Clone()
	}
	newClient.SetHeaderMap(m)
	return newClient
}

// HeaderRaw is a chaining function,
// which sets custom HTTP header using raw string for next request.
func (c *Client) HeaderRaw(headers string) *Client {
	newClient := c
	if c.parent == nil {
		newClient = c.Clone()
	}
	newClient.SetHeaderRaw(headers)
	return newClient
}

func (c *Client) SetHeader(key, value string) *Client {
	c.header[key] = value
	return c
}

// SetHeaderMap sets custom HTTP headers with map.
func (c *Client) SetHeaderMap(m map[string]string) *Client {
	for k, v := range m {
		c.header[k] = v
	}
	return c
}

// SetHeaderRaw sets custom HTTP header using raw string.
func (c *Client) SetHeaderRaw(headers string) *Client {
	for _, line := range convert.SplitAndTrim(headers, "\n") {
		array, _ := gregex.MatchString(`^([\w\-]+):\s*(.+)`, line)
		if len(array) >= 3 {
			c.header[array[1]] = array[2]
		}
	}
	return c
}

type ReadCloser struct {
	index      int    // Current read position.
	content    []byte // Content.
	repeatable bool
}

// NewReadCloser creates and returns a RepeatReadCloser object.
func NewReadCloser(content []byte, repeatable bool) io.ReadCloser {
	return &ReadCloser{
		content:    content,
		repeatable: repeatable,
	}
}

// Read implements the io.ReadCloser interface.
func (b *ReadCloser) Read(p []byte) (n int, err error) {
	n = copy(p, b.content[b.index:])
	b.index += n
	if b.index >= len(b.content) {
		// Make it repeatable reading.
		if b.repeatable {
			b.index = 0
		}
		return n, io.EOF
	}
	return n, nil
}

// Close implements the io.ReadCloser interface.
func (b *ReadCloser) Close() error {
	return nil
}

// clientMiddleware is the plugin for http client request workflow management.
type clientMiddleware struct {
	client       *Client       // http client.
	handlers     []HandlerFunc // mdl handlers.
	handlerIndex int           // current handler index.
	resp         *Response     // save resp.
	err          error         // save err.
}

func (m *clientMiddleware) Next(req *http.Request) (resp *Response, err error) {
	if m.err != nil {
		return m.resp, m.err
	}
	if m.handlerIndex < len(m.handlers) {
		m.handlerIndex++
		m.resp, m.err = m.handlers[m.handlerIndex](m.client, req)
	}
	return m.resp, m.err
}

func (r *Response) ReadAll() []byte {
	// Response might be nil.
	if r == nil || r.Response == nil {
		return []byte{}
	}
	body, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		return nil
	}
	return body
}

func (r *Response) ReadAllString() string {
	return convert.UnsafeBytesToStr(r.ReadAll())
}
