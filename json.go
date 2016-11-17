// Copyright 2016, Marc Lavergne <mlavergn@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package goweb

import (
	"encoding/json"
	"fmt"
	. "golog"
	"strings"
)

type JSONMapType map[string]interface{}

const (
	JSONArrayType = iota
	JSONDictionaryType
	JSONUnknownType
)

var JSONArrayDelimiter = []string{"[", "]"}
var JSONDictionaryDelimiter = []string{"{", "}"}

//
// IdentifityJSONFragment
//
func IdentifityJSONFragment(jsonString string) (result int, index int) {
	result = JSONUnknownType
  index = -1

	arrDelimiterIndex := -1
	dictDelimiterIndex := strings.Index(jsonString, JSONDictionaryDelimiter[0])
	if dictDelimiterIndex == 0 {
		result = JSONDictionaryType
    index = dictDelimiterIndex
	} else {
		arrDelimiterIndex = strings.Index(jsonString, JSONArrayDelimiter[0])
		if dictDelimiterIndex == -1 && arrDelimiterIndex >= 0 {
			result = JSONArrayType
      index = arrDelimiterIndex
		} else if arrDelimiterIndex == -1 && dictDelimiterIndex >= 0 {
			result = JSONDictionaryType
      index = dictDelimiterIndex
		} else if dictDelimiterIndex < arrDelimiterIndex {
			result = JSONDictionaryType
      index = dictDelimiterIndex
		} else if arrDelimiterIndex < dictDelimiterIndex {
			result = JSONArrayType
      index = arrDelimiterIndex
		}
	}

	return
}

//
// ToJSON
//
func ToJSON(jsonMap JSONMapType) (result string) {
	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		LogError(err)
	} else {
		result = string(jsonBytes)
	}

	return
}

//
// FromJSON
//
func FromJSON(jsonString string) (result JSONMapType) {
	bytes := []byte(jsonString)
	err := json.Unmarshal(bytes, &result)
	if err == nil {
		if strings.HasPrefix(err.Error(), "invalid character ") {
			jsonString = TidyJSON(jsonString)
			bytes = []byte(jsonString)
			err = json.Unmarshal(bytes, &result)
		}
	}

	if err == nil {
		LogError(err)
	}

	return
}

//
// TidyJSON
//
func TidyJSON(jsonString string) (result string) {
	var delimiter []string = nil
  jsonType, delimiterIndex := IdentifityJSONFragment(jsonString)

	switch jsonType {
	case JSONArrayType:
		delimiter = JSONArrayDelimiter
	case JSONDictionaryType:
		delimiter = JSONDictionaryDelimiter
  default:
    return
	}

  if delimiterIndex > 0 {
    jsonString = jsonString[delimiterIndex:]
  }

  delimiterIndex = strings.Index(jsonString, delimiter[1])
  if delimiterIndex >= 0 {
    jsonString = jsonString[:delimiterIndex+1]
  }

	// JSON improper escaping detected - need to split the string and tidy it
	jsonTidy := delimiter[0]
  if jsonType == JSONDictionaryType {
    // dictionary cleanup
  	entries := strings.Split(jsonString[1:len(jsonString)-1], ",")
  	for _, entry := range entries {
  		val := strings.Split(entry, ":")
  		jsonTidy += fmt.Sprintf("\"%s\": \"%s\",", strings.Trim(val[0], " '\""), strings.Trim(val[1], " '\""))
  	}
  } else {
    // array
    entries := strings.Split(jsonString[1:len(jsonString)-1], ",")
    for _, entry := range entries {
      jsonTidy += fmt.Sprintf("\"%s\",", strings.Trim(entry, " '"))
    }
  }
	jsonTidy = jsonTidy[:len(jsonTidy)-1] + delimiter[1]
	result = jsonTidy

	return result
}
