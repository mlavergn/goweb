// Copyright 2016, Marc Lavergne <mlavergn@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package goweb

import (
	"encoding/json"
	"errors"
	"fmt"
	. "golog"
	"strconv"
	"strings"
)

type JSONSliceType []interface{}
type JSONMapType map[string]interface{}

const (
	JSONArrayType = iota
	JSONDictionaryType
	JSONUnknownType
)

type _JSONDelimiter []string

var _JSONArrayDelimiter = []string{"[", "]"}
var _JSONDictionaryDelimiter = []string{"{", "}"}

//
// IdentifityJSONFragment
//
func IdentifityJSONFragment(jsonString string) (result int, index int) {
	result = JSONUnknownType
	index = -1

	arrDelimiterIndex := -1
	dictDelimiterIndex := strings.Index(jsonString, _JSONDictionaryDelimiter[0])
	if dictDelimiterIndex == 0 {
		result = JSONDictionaryType
		index = dictDelimiterIndex
	} else {
		arrDelimiterIndex = strings.Index(jsonString, _JSONArrayDelimiter[0])
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
func ToJSON(jsonMap JSONMapType) (result string, err error) {
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
func FromJSON(jsonString string) (result JSONMapType, err error) {
	bytes := []byte(jsonString)
	err = json.Unmarshal(bytes, &result)

	if err != nil {
		if strings.Index(err.Error(), "cannot unmarshal array into") != -1 {
			var resultArr JSONSliceType
			err = json.Unmarshal(bytes, &resultArr)
			if err == nil {
				result = make(JSONMapType)
				result["[]"] = resultArr
			}
		}
	}

	return
}

//
// ExtractJSON isolates and tidys JSON string while attepting to parse the contents into a JSONMap
//
func ExtractJSON(jsonString string, jsonType int) (result JSONMapType, err error) {
	// unmarshall is strict and wants complete JSON structures
	jsonString, jsonType = IsolateJSON(jsonString, jsonType)
	result, err = FromJSON(jsonString)
	if err != nil {
		jsonString = TidyScript(jsonString)
		result, err = FromJSON(jsonString)
		if err != nil {
			jsonString = TidyJSON(jsonString, jsonType)
			result, err = FromJSON(jsonString)
		}
	}

	return
}

//
// IsolateJSON isolates the JSON parsable text based on array or dictionary delimiters
//
func IsolateJSON(jsonString string, jsonTypeIn int) (result string, jsonType int) {
	var delimiter []string = nil
	var delimiterIndex int

	if jsonTypeIn == JSONUnknownType {
		jsonType, delimiterIndex = IdentifityJSONFragment(jsonString)
	} else {
		jsonType = jsonTypeIn
		delimiterIndex = strings.Index(jsonString, _JSONDictionaryDelimiter[0])

	}

	switch jsonType {
	case JSONArrayType:
		delimiter = _JSONArrayDelimiter
	case JSONDictionaryType:
		delimiter = _JSONDictionaryDelimiter
	default:
		return
	}

	if delimiterIndex > 0 {
		result = jsonString[delimiterIndex:]
	}

	delimiterIndex = strings.Index(result, delimiter[1])
	if delimiterIndex >= 0 {
		result = result[:delimiterIndex+1]
	}

	return
}

//
// TidyScript
//
func TidyScript(jsonString string) (result string) {
	// no newlines
	result = strings.Replace(jsonString, "\n", "", -1)
	// no tabs
	result = strings.Replace(result, "\t", "", -1)

	return
}

//
// TidyJSON
//
func TidyJSON(jsonString string, jsonType int) (result string) {
	// JSON improper escaping detected - need to split the string and tidy it
	var jsonDelimiter []string
	if jsonType == JSONDictionaryType {
		// dictionary cleanup
		jsonDelimiter = _JSONDictionaryDelimiter
		entries := strings.Split(jsonString[1:len(jsonString)-1], ",")
		for _, entry := range entries {
			val := strings.Split(entry, ":")
			result += fmt.Sprintf("\"%s\": \"%s\",", strings.Trim(val[0], " '\""), strings.Trim(val[1], " '\""))
		}
	} else {
		// array
		jsonDelimiter = _JSONArrayDelimiter
		entries := strings.Split(jsonString[1:len(jsonString)-1], ",")
		for _, entry := range entries {
			result += fmt.Sprintf("\"%s\",", strings.Trim(entry, " '"))
		}
	}

	// reconstitute the result dropping the last comma
	result = jsonDelimiter[0] + result[:len(result)-1] + jsonDelimiter[1]

	return result
}

func numericConversion(jsonstring string) (result int, err error) {
	if strings.ToLower(jsonstring) == strings.ToUpper(jsonstring) {
		result, err = strconv.Atoi(jsonstring)
		if err != nil {
			result, err = EvaluateEquation(jsonstring)
		}
	} else {
		err = errors.New("string")
	}

	return
}

//
// TidyValues
//
func TidyValues(jsonString string, jsonType int) (result string) {
	var jsonDelimiter []string
	if jsonType == JSONDictionaryType {
		// dictionary cleanup
		jsonDelimiter = _JSONDictionaryDelimiter
		entries := strings.Split(jsonString[1:len(jsonString)-1], ",")
		for _, entry := range entries {
			val := strings.Split(entry, ":")
			ival, err := numericConversion(val[1])
			if err == nil {
				result += fmt.Sprintf("\"%s\": %d,", val[0], ival)
			} else {
				result += fmt.Sprintf("\"%s\": \"%s\",", val[0], val[1])
			}
		}
	} else {
		// array
		jsonDelimiter = _JSONArrayDelimiter
		entries := strings.Split(jsonString[1:len(jsonString)-1], ",")
		for _, entry := range entries {
			ival, err := numericConversion(entry)
			if err == nil {
				result += fmt.Sprintf("%d,", ival)
			} else {
				result += fmt.Sprintf("\"%s\",", entry)
			}
		}
	}

	// reconstitute the result dropping the last comma
	result = jsonDelimiter[0] + result[:len(result)-1] + jsonDelimiter[1]

	return result
}
