// Copyright 2016, Marc Lavergne <mlavergn@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package goweb

import (
  . "golog"
  "testing"
)

func TestTidyJSON(t *testing.T) {
  SetLogLevel(LOG_DEBUG)
  r := TidyJSON("var x = {abc:123, def:\"789\"}")
  LogDebug(len(r))
  if len(r) != 27 {
    t.Errorf("JSON string lenght %d vs expected %d",len(r), 27)
  }
}
