// Copyright 2016, Marc Lavergne <mlavergn@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package goweb

import (
	. "golog"
	"testing"
)

func TestEvaluateEquation(t *testing.T) {
	SetLogLevel(LOG_DEBUG)
	r, err := EvaluateEquation("(3 + 4) * -2 + 10")
	if err != nil {
		t.Errorf("Error %s [%s]", err, r)
	} else {
		if r != -4 {
			t.Errorf("Evaluation %d vs expected %d", r, -4)
		}
	}
}
