// +build !appengine

package goweb

import (
	"net"
	"net/http"
	"time"
)

//
// GetClient takes an unused http.Request so we can support GAE
//
func GetClient(r *http.Request) (client *http.Client) {
	client = &http.Client{
		Timeout: time.Second * 30,
	}

	return
}

//
// Assumes a local HTTP proxy is anything listening locally on port 8080
//
func DetectProxy() (result bool) {
	// if GAE is running, port 8000 will be active, so deny proxying for now
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		conn, err = net.Dial("tcp", "127.0.0.1:8080")
		if err != nil {
			result = false
		} else {
			conn.Close()
			result = true
		}
	} else {
		conn.Close()
		result = false
	}

	return
}
