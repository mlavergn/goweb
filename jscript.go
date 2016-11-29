// Copyright 2016, Marc Lavergne <mlavergn@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package goweb

import (
	. "golog"
	"regexp"
)

type JScript struct {
}

//
// Constructor
//
func NewJScript() *JScript {
	return &JScript{}
}

//
//
//
func (self *JScript) ParseRedirect(d *DOM) string {
	var result string

	scripts := d.Find("script", nil)
	if len(scripts) > 0 {
		LogDebug("SCRIPTs found")

		re, _ := regexp.Compile("document.location\\s?=\\s?['\"](.+)[\"'];")

		for _, script := range scripts {
			match := re.FindStringSubmatch(script.Text())
			if len(match) > 1 {
				result = match[1]
			}
		}
		if len(result) > 0 {
			LogDebug("Script redirect detected: " + result)
		} else {
			LogDebug("No script redirect detected")
		}
	} else {
		LogDebug("META not found")
	}

	return result
}
