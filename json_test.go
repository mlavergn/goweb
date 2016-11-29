// Copyright 2016, Marc Lavergne <mlavergn@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package goweb

import (
	. "golog"
	"testing"
)

func TestFromJSON(t *testing.T) {
	SetLogLevel(LOG_DEBUG)
	r, err := FromJSON("[11, 22, 33]")
	if err != nil {
		t.Errorf("Error %s [%s]", err, r)
	} else {
		if r["[]"] == nil {
			t.Errorf("JSON length %d vs expected %d [%s]", len(r), 3, r)
		}
	}
}

func TestIsolateJSON(t *testing.T) {
	SetLogLevel(LOG_DEBUG)
	r, _ := IsolateJSON("var x = {abc:123, def:\"789\"};\nvar y = [1, 2, 3+4];", JSONDictionaryType)
	if len(r) != 20 {
		t.Errorf("JSON string length %d vs expected %d [%s]", len(r), 20, r)
	}
}

func TestExtractJSON(t *testing.T) {
	SetLogLevel(LOG_DEBUG)
	r, err := ExtractJSON("var x = {abc:123, def:\"789\"}", JSONDictionaryType)
	if err != nil {
		t.Errorf("Error %s [%s]", err, r)
	} else {
		if len(r) != 2 {
			t.Errorf("JSON map length %d vs expected %d [%s]", len(r), 2, r)
		}
	}
}

func TestTidyValues(t *testing.T) {
	SetLogLevel(LOG_DEBUG)
	r := TidyValues("[100, 300 + -50]", JSONArrayType)
	if len(r) != 9 {
		t.Errorf("JSON array length %d vs expected %d [%s]", len(r), 9, r)
	}
}
