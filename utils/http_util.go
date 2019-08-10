package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	nurl "net/url"
	"strings"
	"time"
)

type Http_Util struct {
	client *http.Client
}

func (c *Http_Util) GetJSON(url string, queryParams map[string]string, responseData interface{}, options *RequestOptions) error {

	u, err := nurl.Parse(url)
	if err != nil {
		return err
	}
	params := u.Query()
	if queryParams != nil {
		for key, value := range queryParams {
			params.Set(key, value)
		}
	}
	u.RawQuery = params.Encode()

	return c.RequestJSON(http.MethodGet, u.String(), nil, responseData, options)
}

func (c *Http_Util) PostJSON(url string, formParams map[string]interface{}, responseData interface{}, options *RequestOptions) error {
	return c.RequestJSON(http.MethodPost, url, formParams, responseData, options)
}

type httpRequestEncodeType int

const (
	URLEncoded  httpRequestEncodeType = 1
	JSONEncoded httpRequestEncodeType = 2
)

func (c *Http_Util) RequestJSON(method string, url string, formParams map[string]interface{}, responseData interface{}, options *RequestOptions) error {

	var byteBuff *bytes.Buffer
	var err error

	encodeType := JSONEncoded

	if options != nil && options.Headers != nil {
		if contentType := options.Headers["Content-Type"]; strings.HasPrefix(contentType,
			"application/x-www-form-urlencoded") {
			encodeType = URLEncoded
		} else {
			encodeType = JSONEncoded
		}
	}

	if formParams != nil {
		if encodeType == URLEncoded {
			urlValues := nurl.Values{}

			for key, value := range formParams {
				urlValues.Add(key, fmt.Sprintf("%v", value))
			}
			byteBuff = bytes.NewBuffer([]byte(urlValues.Encode()))
		} else if encodeType == JSONEncoded {
			byteBuff, err = JSONUtil.MapToByteBuffer(formParams)
			if err != nil {
				return err
			}
		}
	}

	var req *http.Request

	if byteBuff == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {

		req, err = http.NewRequest(method, url, byteBuff)
	}
	if err != nil {
		return err
	}

	if options != nil && options.Headers != nil {
		for key, value := range options.Headers {
			req.Header.Add(key, value)
		}
	}

	resp, err := c.client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var responseBuff []byte
	responseBuff, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.NewDecoder(bytes.NewBuffer(responseBuff)).Decode(responseData)

	if err != nil {
		log.Println(fmt.Sprintf("Error:%s", err.Error()))
		log.Println(fmt.Sprintf("url:%s\ncontent:%s", url, string(responseBuff)))
	}
	return err
}

func (c *Http_Util) SetProxy(host string, port int) {
	proxyURL, err := url.Parse(fmt.Sprintf("http://%s:%d", host, port))
	if err != nil {
		// log
	}
	transport := c.client.Transport.(*http.Transport)
	transport.Proxy = http.ProxyURL(proxyURL)
}

func (c *Http_Util) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

func (c *Http_Util) UseCookieJar(use bool) {
	if use {
		if c.client.Jar == nil {
			jar, _ := cookiejar.New(nil)
			c.client.Jar = jar
		}
	} else {
		c.client.Jar = nil
	}
}

func NewClient() *http.Client {

	transport := &http.Transport{}

	client := &http.Client{Transport: transport, Timeout: time.Duration(90) * time.Second}

	return client
}

var HttpUtil *Http_Util

const (
	httpTimeoutSecondSettingKey = "http_timeout_seoncd"
	httpProxySettingKey         = "http_proxy"
	httpUseCookieJarSettingKey  = "http_use_cookie_jar"

	httpDebugErrorJSON = "http_debug_error_json"
)

func init() {
	HttpUtil = &Http_Util{
		client: NewClient(),
	}
}

type RequestOptions struct {
	Headers       map[string]string
	ContentCipher []byte
}
