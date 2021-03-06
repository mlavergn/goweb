// Copyright 2016, Marc Lavergne <mlavergn@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package goweb

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	. "golog"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

const (
	HTTP_GET  = "GET"
	HTTP_POST = "POST"
)

const (
	CONTENT_TYPE_NONE       = ""
	CONTENT_TYPE_FORM       = "application/x-www-form-urlencoded"
	CONTENT_TYPE_JSON       = "application/json"
	CONTENT_TYPE_FORM_MULTI = "multipart/form-data"
)

var HTTP_USER_AGENT = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12) AppleWebKit/602.1.50 (KHTML, like Gecko) Version/10.0 Safari/602.1.50",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.75.14 (KHTML, like Gecko) Version/7.0.3 Safari/7046A194A",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36",
}

type HTTP struct {
	Method      string
	URL         *url.URL
	ProxyURL    *url.URL
	cookieJar   http.CookieJar
	req         *http.Request
	resp        *http.Response
	RawContents []byte

	gaeRequest *http.Request
}

//
// Constructor
//
func NewHTTP() *HTTP {
	id := &HTTP{}
	id.cookieJar, _ = cookiejar.New(nil)

	return id
}

func (id *HTTP) SetGAERequest(req *http.Request) {
	id.gaeRequest = req
}

func (id *HTTP) Contents() string {
	return string(id.RawContents)
}

//
// Tidy the URL such that it is minimally valid
//
func (id *HTTP) tidyURL(urlString string) (err error) {
	LogDebugf("tidyURL input: %s", urlString)
	url, err := url.Parse(urlString)

	if id.URL == nil {
		if len(url.Scheme) == 0 {
			url.Scheme = "http"
		}

		id.URL = url
	} else {
		if len(url.Scheme) == 0 {
			url.Scheme = id.URL.Scheme
		}

		if len(url.Host) == 0 {
			url.Host = id.URL.Host
		}

		if !strings.HasPrefix(url.Path, "/") {
			components := strings.Split(id.URL.Path, "/")
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

		LogDebugf("tidyURL output: %s", url)
		id.URL = url
	}

	return err
}

func (id *HTTP) URLString() (urlString string) {
	urlString = fmt.Sprintf("%s", id.URL)

	return urlString
}

//
// Fetch: POST request
//
func (id *HTTP) Post(urlString string, args map[string]string) (result string) {
	content := _formatArgs(args)
	result = id.PostContent(urlString, CONTENT_TYPE_NONE, content)

	return result
}

func (id *HTTP) PostString(urlString string, contentType string, contentString string) (result string) {
	content := bytes.NewBuffer([]byte(contentString))
	result = id.PostContent(urlString, contentType, content)

	return result
}

func (id *HTTP) PostData(urlString string, contentType string, contentBytes []byte) (result string) {
	content := bytes.NewBuffer(contentBytes)
	result = id.PostContent(urlString, contentType, content)

	return result
}

//
// Fetch: POST request
//
func (id *HTTP) PostContent(urlString string, contentType string, content *bytes.Buffer) (result string) {
	id.Method = HTTP_POST
	id.tidyURL(urlString)

	result = id.prepareAndExecuteRequest(contentType, content)
	LogDebugf("POST status %d", id.Status())

	return result
}

//
// Fetch: GET request
//
func (id *HTTP) Get(urlString string) (result string) {
	return id.GetQuery(urlString, nil)
}

func (id *HTTP) GetQuery(urlString string, args map[string]string) (result string) {
	id.Method = HTTP_GET
	id.tidyURL(urlString)

	if args != nil {
		qry := url.Values{}
		for key, val := range args {
			qry.Add(key, val)

		}
		id.URL.RawQuery = qry.Encode()
	}

	result = id.prepareAndExecuteRequest(CONTENT_TYPE_NONE, nil)
	LogDebugf("GET status %d", id.Status())

	return result
}

//
//
//
func _formatArgs(args map[string]string) (content *bytes.Buffer) {
	if len(args) > 0 {
		argString := ""
		for key, val := range args {
			argString += key + "=" + val + "&"
		}
		content = bytes.NewBuffer([]byte(argString))
	} else {
		content = bytes.NewBuffer([]byte(""))
	}

	return content
}

//
// Fetch: Prepare and execute HTTP request
// NOTE: This is the work horse, all requests filter through here
//
func (id *HTTP) prepareAndExecuteRequest(contentType string, content *bytes.Buffer) string {
	LogDebug(id.Method + ": " + id.URLString())

	RedirectAttemptedError := errors.New("redirect")

	client := GetClient(id.gaeRequest)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return RedirectAttemptedError
	}
	client.Jar = id.cookieJar

	// if we're proxying, we're going to disable the TLS cert verification
	if DetectProxy() {
		id.ProxyURL, _ = url.Parse("http://127.0.0.1:8080")
		client.Transport = &http.Transport{
			Proxy:           http.ProxyURL(id.ProxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	// ensure we don't pass a nil ByteBuffer into http.NewRequest
	if content == nil {
		content = bytes.NewBuffer([]byte(""))
	}

	id.req, _ = http.NewRequest(id.Method, id.URLString(), content)

	if len(id.URL.Host) > 0 {
		id.setRequestHeader("Host", id.URL.Host)
	}

	if len(contentType) > 0 {
		id.setRequestHeader("Content-Type", contentType)
	}

	id.setRequestHeader("Connection", "close")

	// randomize the user agent
	uaStr := HTTP_USER_AGENT[rand.Intn(len(HTTP_USER_AGENT))]
	id.setRequestHeader("User-Agent", uaStr)

	// id.setRequestHeader("Referer", referrer)
	var err error
	id.resp, err = client.Do(id.req)

	if err != nil && !strings.HasSuffix(err.Error(), " redirect") {
		LogError(err)
		return ""
	} else {
		bytes, _ := ioutil.ReadAll(id.resp.Body)
		id.RawContents = bytes
	}
	defer id.resp.Body.Close()

	// at this point we have the request and response, save a record if configured
	output := "<!--\nMethod: " + id.Method + "\nURL: " + id.URLString() + "\nStatus: " + strconv.Itoa(id.Status()) + "\n-->\n\n" + id.Contents()
	LogDumpFile("goweb", output)

	// handle redirects
	id.handleRedirection()

	return id.Contents()
}

//
// Handler: redirections
//
func (id *HTTP) handleRedirection() string {
	var result string

	switch id.Status() {
	case 200:
		// OK
		if id.isHTML() {
			LogDebug("HTML detected")
			h := NewHTML()
			s := NewDOM()
			s.SetContents(id.Contents())
			url := h.ParseRedirect(s)
			if len(url) > 0 {
				id.Get(url)
				result = id.handleRedirection()
			}
		} else if id.isJSON() {
			LogDebug("JSON detected")
			_, err := id.JSON()
			if err == nil {
				result = id.Contents()
			}
		} else {
			LogDebug("Unhandled content type detected: " + id.ContentType())
			result = id.Contents()
		}
		break
	case 302:
		// MOVED
		url := id.Location()
		id.Get(url)
		result = id.handleRedirection()
		break
	default:
		LogWarn("Unhandled status")
	}

	return result
}

//
// Contents: Determine if the contents are JSON
//
func (id *HTTP) isHTML() (result bool) {
	result = false
	if strings.HasPrefix(id.ContentType(), "text/html") {
		result = true
	}

	return result
}

//
// Contents: Determine if the contents are JSON
//
func (id *HTTP) isJSON() (result bool) {
	result = false
	if strings.HasPrefix(id.ContentType(), "application/json") {
		result = true
	}

	return result
}

//
// Contents: Determine if the contents are XML
//
func (id *HTTP) isXML() (result bool) {
	result = false
	if strings.HasPrefix(id.ContentType(), "text/xml") {
		result = true
	}

	return result
}

//
// Contents: Determine if the contents are XML
//
func (id *HTTP) isImage() (result bool) {
	result = false
	if strings.HasPrefix(id.ContentType(), "image/") {
		result = true
	}

	return result
}

//
// Contents: Marshall contents from JSON to a map if possible
//
func (id *HTTP) JSON() (result map[string]interface{}, err error) {
	err = json.Unmarshal(id.RawContents, &result)

	return
}

//
// Header: Extract Content-Type
//
func (id *HTTP) ContentType() string {
	return id.GetResponseHeader("Content-Type")
}

//
// Header: Extract Hosts
//
func (id *HTTP) Host() string {
	return id.GetResponseHeader("Host")
}

//
// Header: Extract Status
//
func (id *HTTP) Status() int {
	if id.resp == nil {
		LogWarn("Response is not available")
		// return 499 - Client Closed Request
		return 499
	}

	return id.resp.StatusCode
}

//
// Header: Extract Location
//
func (id *HTTP) Location() string {
	return id.GetResponseHeader("Location")
}

//
// Header: Extract Referer
//
func (id *HTTP) Referer() string {
	return id.GetResponseHeader("Referer")
}

//
// Header: Extracts a header by key from the response
//
func (id *HTTP) GetResponseHeader(key string) string {
	if id.resp == nil {
		LogWarn("Response is not available")
		return ""
	}

	var result string = ""
	varr := id.resp.Header[key]
	if len(varr) > 0 {
		result = varr[0]
	}
	return result
}

//
// Header: Sets a header by key in the request
//
func (id *HTTP) setRequestHeader(key string, value string) {
	if id.req == nil {
		LogWarn("Request is not ready")
		return
	}

	id.req.Header.Add(key, value)
}
