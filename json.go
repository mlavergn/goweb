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
func IdentifityJSONFragment(jsonString string) (result int) {
	result = JSONUnknownType
	arrFrontPriority := -1
	dictFrontPriority := strings.Index(jsonString, JSONDictionaryDelimiter[0])
	if dictFrontPriority == 0 {
		result = JSONDictionaryType
	} else {
		arrFrontPriority = strings.Index(jsonString, JSONArrayDelimiter[0])
		if dictFrontPriority == -1 && arrFrontPriority >= 0 {
			result = JSONArrayType
		} else if arrFrontPriority == -1 && dictFrontPriority >= 0 {
			result = JSONDictionaryType
		} else if dictFrontPriority < arrFrontPriority {
			result = JSONDictionaryType
		} else if arrFrontPriority < dictFrontPriority {
			result = JSONArrayType
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
	switch IdentifityJSONFragment(jsonString) {
	case JSONArrayType:
		delimiter = JSONArrayDelimiter
	case JSONDictionaryType:
		delimiter = JSONDictionaryDelimiter
	}

	// JSON improper escaping detected - need to split the string and tidy it
	LogDebug("Tidy JSON")
	jsonTidy := delimiter[0]
	entries := strings.Split(jsonString[1:len(jsonString)-1], ",")
	for _, entry := range entries {
		val := strings.Split(entry, ":")
		jsonTidy += fmt.Sprintf("\"%s\": \"%s\",", strings.Trim(val[0], " '"), strings.Trim(val[1], " '"))
	}
	jsonTidy = jsonTidy[:len(jsonTidy)-1] + delimiter[1]
	result = jsonTidy

	return result
}
