package main

import (
	"fmt"
)

var (
	ErrNotSuitableToken = "token %q is not suitable, expected token type is %q\n"
	ErrNoTokenExists    = "there is no token, but %q need to exist\n"
)

type CompilationEngine struct {
	input []*Token
	pos   int
}

func NewCompilationEngine(tokens []*Token) *CompilationEngine {
	if len(tokens) == 0 {
		return &CompilationEngine{}
	}
	return &CompilationEngine{
		input: tokens,
		pos:   0,
	}
}

func (c *CompilationEngine) Parse() (*Token, error) {
	output, err := c.compileClass()
	if err != nil {
		return &Token{}, err
	}
	return output, nil
}

func (c *CompilationEngine) compileClass() (*Token, error) {
	// ["class"] must be keyword "class"
	t, ok := c.current()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "class")
	}
	if t.Value != "class" || t.Type != keyword {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, class)
	}
	classToken := &Token{
		Type:     class,
		Children: make([]*Token, 0),
	}
	classToken.Children = append(classToken.Children, t)

	// [className] must be identifier
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "className")
	}
	if t.Type != identifier {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, class)
	}
	classToken.Children = append(classToken.Children, t)

	// ["{"] must be simbol: "{"
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "{")
	}
	if t.Value != "{" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	classToken.Children = append(classToken.Children, t)

LOOP:
	for t, ok := c.next(); ok; t, ok = c.next() {
		switch t.Value {
		case "static", "field":
			if t.Type != keyword {
				return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, keyword)
			}
			c.rewind()
			classVar, err := c.compileClassVarDec()
			if err != nil {
				return nil, err
			}
			classToken.Children = append(classToken.Children, classVar)
		case "constructor", "function", "method":
			if t.Type != keyword {
				return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, keyword)
			}
			c.rewind()
			subrDec, err := c.compileSubroutineDec()
			if err != nil {
				return nil, err
			}
			classToken.Children = append(classToken.Children, subrDec)
		case "}":
			if t.Type != symbol {
				return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
			}
			classToken.Children = append(classToken.Children, t)
			break LOOP
		default:
			return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, keyword)
		}
	}
	fmt.Printf("complete parse all tokens. inout tokens: %v, parsed tokens: %v\n", len(c.input), c.pos)
	return classToken, nil
}

func (c *CompilationEngine) compileClassVarDec() (*Token, error) {
	// ["static" | "field"] must be keyword "static" or "field"
	t, ok := c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "static | field")
	}
	if (t.Value != "static" && t.Value != "field") || t.Type != keyword {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, "static, field")
	}
	classVar := &Token{
		Type:     classVarDec,
		Children: make([]*Token, 0),
	}
	classVar.Children = append(classVar.Children, t)

	// [type] must be keyword "int" or "char" or "boolean" or className
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "type")
	}
	if !isType(t) {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, "type")
	}
	classVar.Children = append(classVar.Children, t)

	// [varName] must be identifier
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "varName")
	}
	if t.Type != identifier {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, identifier)
	}
	classVar.Children = append(classVar.Children, t)

	// [", varName" | ";"] must be keyword "," or ";" or varName
LOOP:
	for {
		t, ok := c.next()
		if !ok {
			return nil, fmt.Errorf(ErrNoTokenExists, ", | ;")
		}
		switch t.Value {
		case ",":
			if t.Type != symbol {
				return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
			}
			classVar.Children = append(classVar.Children, t)
			t2, ok := c.next()
			if !ok {
				return nil, fmt.Errorf(ErrNoTokenExists, "varName")
			}
			if t2.Type != identifier {
				return nil, fmt.Errorf(ErrNotSuitableToken, t2.Value, identifier)
			}
			classVar.Children = append(classVar.Children, t2)
		case ";":
			if t.Type != symbol {
				return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
			}
			classVar.Children = append(classVar.Children, t)
			break LOOP
		}
	}

	return classVar, nil
}

func (c *CompilationEngine) compileSubroutineDec() (*Token, error) {
	// ["constructor" | "function" | "method"] must be keyword "constructor" or "function" or "method"
	t, ok := c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "constructor | function | method")
	}
	if (t.Value != "constructor" && t.Value != "function" && t.Value != "method") || t.Type != keyword {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, "static, field")
	}
	subrDec := &Token{
		Type:     subroutineDec,
		Children: make([]*Token, 0),
	}
	subrDec.Children = append(subrDec.Children, t)

	// ["void" | type] must be keyword "int" or "char" or "boolean" or className
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "constructor | function | method")
	}
	if t.Value == "void" {
		if t.Type != keyword {
			return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, keyword)
		}
		subrDec.Children = append(subrDec.Children, t)
	} else if isType(t) {
		subrDec.Children = append(subrDec.Children, t)
	} else {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, "type | void")
	}

	// [subroutineName] must be identifier
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "subroutineName")
	}
	if t.Type != identifier {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, identifier)
	}
	subrDec.Children = append(subrDec.Children, t)

	// ["("] must be simbol: "("
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "(")
	}
	if t.Value != "(" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	subrDec.Children = append(subrDec.Children, t)

	// [parameterList] if exists
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "type | )")
	}
	if isType(t) {
		paramList, err := c.compileParameterList()
		if err != nil {
			return nil, err
		}
		subrDec.Children = append(subrDec.Children, paramList)
		t, ok = c.next()
		if !ok {
			return nil, fmt.Errorf(ErrNoTokenExists, ")")
		}
	}

	// [")"] must be simbol: ")"
	if t.Value != ")" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	subrDec.Children = append(subrDec.Children, t)

	// [subroutineBody]
	subrBody, err := c.compileSubroutineBody()
	if err != nil {
		return nil, err
	}
	subrDec.Children = append(subrDec.Children, subrBody)

	return subrDec, nil
}

func isType(t *Token) bool {
	switch t.Value {
	case "int", "char", "boolean":
		if t.Type != keyword {
			return false
		}
	default:
		// className
		if t.Type != identifier {
			return false
		}
	}
	return true
}

func (c *CompilationEngine) compileParameterList() (*Token, error) {
	// [type] must be keyword or identity
	t, ok := c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "type")
	}
	if !isType(t) {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, "type")
	}
	paramList := &Token{
		Type:     parameterList,
		Children: make([]*Token, 0),
	}
	paramList.Children = append(paramList.Children, t)

	// [varName] must be identifier
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "varName")
	}
	if t.Type != identifier {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, identifier)
	}
	paramList.Children = append(paramList.Children, t)

	// ["," varName | ")"] must be keyword "," varName or ")"
LOOP:
	for {
		t, ok := c.next()
		if !ok {
			return nil, fmt.Errorf(ErrNoTokenExists, ", | )")
		}
		switch t.Value {
		case ",":
			if t.Type != symbol {
				return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
			}
			paramList.Children = append(paramList.Children, t)
			t2, ok := c.next()
			if !ok {
				return nil, fmt.Errorf(ErrNoTokenExists, "varName")
			}
			if t2.Type != identifier {
				return nil, fmt.Errorf(ErrNotSuitableToken, t2.Value, identifier)
			}
			paramList.Children = append(paramList.Children, t2)
		case ")":
			// parameterList is not contained teminal ")"
			c.rewind()
			break LOOP
		}
	}
	return paramList, nil
}

func (c *CompilationEngine) compileSubroutineBody() (*Token, error) {
	// ["{"] must be keyword "{"
	t, ok := c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "{")
	}
	if t.Value != "{" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	subrBody := &Token{
		Type:     subroutineBody,
		Children: make([]*Token, 0),
	}
	subrBody.Children = append(subrBody.Children, t)

	// varDec
	varDecTokens, err := c.compileVarDec()
	if err != nil {
		return nil, err
	}
	subrBody.Children = append(subrBody.Children, varDecTokens)

	// TODO statements

	// ["}"] must be keyword "{"
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "}")
	}
	if t.Value != "}" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	return subrBody, nil
}

func (c *CompilationEngine) compileVarDec() (*Token, error) {
	// ["var"] must be keyword "var"
	t, ok := c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "var")
	}
	if t.Value != "var" || t.Type != keyword {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, keyword)
	}
	varDecTokens := &Token{
		Type:     varDec,
		Children: make([]*Token, 0),
	}
	varDecTokens.Children = append(varDecTokens.Children, t)

	// [type] must be keyword or identity
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "type")
	}
	if !isType(t) {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, "type")
	}
	varDecTokens.Children = append(varDecTokens.Children, t)

	// [varName] must be identifier
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "varName")
	}
	if t.Type != identifier {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, identifier)
	}
	varDecTokens.Children = append(varDecTokens.Children, t)

	// [", varName" | ";"] must be keyword "," or ";" or varName
LOOP:
	for {
		t, ok := c.next()
		if !ok {
			return nil, fmt.Errorf(ErrNoTokenExists, ", | ;")
		}
		switch t.Value {
		case ",":
			if t.Type != symbol {
				return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
			}
			varDecTokens.Children = append(varDecTokens.Children, t)
			t2, ok := c.next()
			if !ok {
				return nil, fmt.Errorf(ErrNoTokenExists, "varName")
			}
			if t2.Type != identifier {
				return nil, fmt.Errorf(ErrNotSuitableToken, t2.Value, identifier)
			}
			varDecTokens.Children = append(varDecTokens.Children, t2)
		case ";":
			if t.Type != symbol {
				return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
			}
			varDecTokens.Children = append(varDecTokens.Children, t)
			break LOOP
		}
	}

	return varDecTokens, nil
}

func (c *CompilationEngine) current() (*Token, bool) {
	if len(c.input) == 0 {
		fmt.Println("there are no tokens")
		return nil, false
	}
	return c.input[c.pos], true
}

func (c *CompilationEngine) next() (*Token, bool) {
	if c.hasNext() {
		fmt.Println("there is no next token")
		return nil, false
	}
	c.pos++
	t := c.input[c.pos]
	return t, true
}

func (c *CompilationEngine) hasNext() bool { return c.pos+1 > len(c.input) }

func (c *CompilationEngine) rewind() {
	if c.pos > 0 {
		c.pos--
	}
}
