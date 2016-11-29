// Copyright 2016, Marc Lavergne <mlavergn@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package goweb

import (
	"go/ast"
	"go/parser"
	"go/token"
	. "golog"
	"regexp"
	"strconv"
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

func EvaluateEquation(javascript string) (result int, err error) {
	var stxTree ast.Expr
	stxTree, err = parser.ParseExpr(javascript)
	if err == nil {
		result = eval(stxTree)
	} else {
		LogError(err)
	}

	return
}

func eval(expr ast.Expr) (result int) {
	switch expr := expr.(type) {
	case *ast.ParenExpr:
		result = eval(expr.X)
	case *ast.BinaryExpr:
		result = evalBinaryExpr(expr)
	case *ast.UnaryExpr:
		result = evalUnaryExpr(expr)
	case *ast.BasicLit:
		switch expr.Kind {
		case token.INT:
			result, _ = strconv.Atoi(expr.Value)
		}
	}

	return
}

func evalBinaryExpr(expr *ast.BinaryExpr) (result int) {
	x := eval(expr.X)
	y := eval(expr.Y)

	switch expr.Op {
	case token.MUL:
		result = x * y
	case token.QUO:
		result = x / y
	case token.ADD:
		result = x + y
	case token.SUB:
		result = x - y
	}

	return
}

func evalUnaryExpr(expr *ast.UnaryExpr) (result int) {
	x := eval(expr.X)

	switch expr.Op {
	case token.ADD:
		result = x
	case token.SUB:
		result = x * -1
	}

	return
}
