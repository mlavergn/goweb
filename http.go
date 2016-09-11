// Copyright 2016, Marc Lavergne <mlavergn@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mnet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"godom"
	. "golog"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

type Enum string

const (
	HTTP_GET  = "GET"
	HTTP_POST = "POST"
)

type HTTP struct {
	Method    string
	URL       *url.URL
	ProxyURL  *url.URL
	cookieJar http.CookieJar
	req       *http.Request
	resp      *http.Response
	Contents  string
}

//
// Constructor
//
func NewHTTP() *HTTP {
	self := &HTTP{}
	self.cookieJar, _ = cookiejar.New(nil)

	return self
}

//
// Tidy the URL such that it is minimally valid
//
func (self *HTTP) _tidyURL(urlString string) (err error) {
	url, err := url.Parse(urlString)

	if self.URL == nil {
		self.URL = url
	} else {
		if len(url.Scheme) == 0 {
			url.Scheme = self.URL.Scheme
			url.Host = self.URL.Host
		}

		if !strings.HasPrefix(url.Path, "/") {
			components := strings.Split(self.URL.Path, "/")
			maxComponent := len(components)
			path := ""
			for i, component := range components {
				if i < maxComponent-1 {
					if len(component) > 0 {
						path += "/" + component
					}
				}
			}

			path += "/" + url.Path
			url.Path = path
		}

		self.URL = url
	}

	return err
}

func (self *HTTP) URLString() (urlString string) {
	urlString = fmt.Sprintf("%s", self.URL)

	return urlString
}

//
// Fetch: POST request
//
func (self *HTTP) Post(urlString string, args map[string]string) string {
	self.Method = HTTP_POST
	self._tidyURL(urlString)

	result := self.prepareAndExecuteRequest(args)
	LogDebugf("\t%d", self.Status())
	LogDebug("\t" + self.ContentType())

	LogDebug(self.Location())

	return result
}

//
// Fetch: GET request
//
func (self *HTTP) Get(urlString string) string {
	self.Method = HTTP_GET
	self._tidyURL(urlString)

	result := self.prepareAndExecuteRequest(nil)
	LogDebugf("\t%d", self.Status())
	LogDebug("\t" + self.ContentType())

	LogDebug(self.Location())

	return result
}

//
// Fetch: Prepare and execute HTTP request
// NOTE: This is the work horse, all requests filter through here
//
func (self *HTTP) prepareAndExecuteRequest(args map[string]string) string {
	LogDebug(self.Method + ": " + self.URLString())

	RedirectAttemptedError := errors.New("redirect")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// return http.ErrUseLastResponse -- Go 1.7+ only
			return RedirectAttemptedError
		},
		Jar: self.cookieJar,
	}

	if self.detectProxy() {
		self.ProxyURL, _ = url.Parse("http://127.0.0.1:8080")
		client.Transport = &http.Transport{Proxy: http.ProxyURL(self.ProxyURL)}
	}

	var argBytes *bytes.Buffer = nil
	if len(args) > 0 {
		argString := ""
		for key := range args {
			argString += key + "=" + "&"
		}
		argBytes = bytes.NewBuffer([]byte(argString))
	} else {
		argBytes = bytes.NewBuffer([]byte(""))
	}

	self.req, _ = http.NewRequest(self.Method, self.URLString(), argBytes)

	self.req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12) AppleWebKit/602.1.50 (KHTML, like Gecko) Version/10.0 Safari/602.1.50")
	// self.req.Header.Add("Referer", referrer)
	var err error
	self.resp, err = client.Do(self.req)

	if err != nil && !strings.HasSuffix(err.Error(), " redirect") {
		LogFatal(err)
	} else {
		bytes, _ := ioutil.ReadAll(self.resp.Body)
		self.Contents = string(bytes)
	}
	// defer s.resp.Body.Close()

	// at this point we have the request and response, save a record if configured
	output := "<!--\nMethod: " + self.Method + "\nURL: " + self.URLString() + "\nStatus: " + strconv.Itoa(self.Status()) + "\n-->\n\n" + self.Contents
	LogDumpFile("mnet", output)

	// handle redirects
	self.handleRedirection()

	return self.Contents
}

//
// Handler: proxies
// Assumes a local HTTP proxy is anything listening locally on port 8080
//
func (self *HTTP) detectProxy() bool {
	result := true

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		result = false
	} else {
		conn.Close()
	}

	return result
}

//
// Handler: redirections
//
func (self *HTTP) handleRedirection() string {
	var result string

	switch self.Status() {
	case 200:
		// OK
		if self.isHTML() {
			LogDebug("HTML detected")
			h := NewHTML()
			s := godom.NewDOM()
			s.SetContents(self.Contents)
			url := h.ParseRedirect(s)
			if len(url) > 0 {
				self.Get(url)
				result = self.handleRedirection()
			}
		} else if self.isJSON() {
			LogDebug("JSON detected")
			data := self.JSON()
			if len(data) > 0 {
				result = self.Contents
			}
		} else {
			LogDebug("Unhandled content type detected: " + self.ContentType())
			result = self.Contents
		}
		break
	case 302:
		// MOVED
		url := self.Location()
		self.Get(url)
		result = self.handleRedirection()
		break
	default:
		LogWarn("Unhandled status")
	}

	return result
}

//
// Contents: Determine if the contents are JSON
//
func (self *HTTP) isHTML() bool {
	var result bool = false
	if strings.HasPrefix(self.ContentType(), "text/html") {
		result = true
	}

	return result
}

//
// Contents: Determine if the contents are JSON
//
func (self *HTTP) isJSON() bool {
	var result bool = false
	if strings.HasPrefix(self.ContentType(), "application/json") {
		result = true
	}

	return result
}

//
// Contents: Marshall contents from JSON to a map if possible
//
func (self *HTTP) JSON() map[string]interface{} {
	var result map[string]interface{}

	if self.isJSON() {
		bytes := []byte(self.Contents)
		json.Unmarshal(bytes, &result)
	}
	return result
}

//
// Header: Extract Content-Type
//
func (self *HTTP) ContentType() string {
	return self.getHeaderValue("Content-Type")
}

//
// Header: Extract Hosts
//
func (self *HTTP) Host() string {
	return self.getHeaderValue("Host")
}

//
// Header: Extract Status
//
func (self *HTTP) Status() int {
	return self.resp.StatusCode
}

//
// Header: Extract Location
//
func (self *HTTP) Location() string {
	return self.getHeaderValue("Location")
}

//
// Header: Extracts a header by key from the response
//
func (self *HTTP) getHeaderValue(key string) string {
	var result string = ""
	varr := self.resp.Header[key]
	if len(varr) > 0 {
		result = varr[0]
	}
	return result
}
