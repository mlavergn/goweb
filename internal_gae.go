// +build appengine

package goweb

import (
	"appengine"
	"appengine/urlfetch"
	"net/http"
)

//
// GetClient GAE uses port 8080, so deny proxing for now
//
func GetClient(r *http.Request) (client *http.Client) {
	ctx := appengine.NewContext(r)
	client = urlfetch.Client(ctx)

	return
}

//
// DetectProxy GAE uses port 8080, so deny proxing for now
//
func DetectProxy() (result bool) {
	result = false

	return
}
