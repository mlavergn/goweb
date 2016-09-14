// Copyright 2016, Marc Lavergne <mlavergn@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package goweb

import (
	"godom"
	. "golog"
	"strings"
)

type HTML struct {
}

//
// Constructor
//
func NewHTML() *HTML {
	return &HTML{}
}

func (self *HTML) ParseRedirect(d *godom.DOM) string {
	var result string

	meta := d.Find("meta", nil)
	if len(meta) > 0 {
		LogDebug("META found")
		value := meta[0].Attr("content")
		token := "http"
		idx := strings.Index(value, token)
		if idx > -1 {
			result = value[idx:]
			LogDebug("META URL detected: " + result)
		} else {
			LogDebug("META no URL detected")
		}
	} else {
		LogDebug("META not found")
	}

	js := NewJScript()
	result = js.ParseRedirect(d)

	return result
}
