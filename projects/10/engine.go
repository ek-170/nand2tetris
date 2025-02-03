package main

import (
	"fmt"
	"unicode"
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
		return &Token{}, fmt.Errorf("pos is %v: %w", c.pos, err)
	}
	return output, nil
}

// ------------- program structures -------------

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

	// [parameterList]
	paramList, err := c.compileParameterList()
	if err != nil {
		return nil, err
	}
	subrDec.Children = append(subrDec.Children, paramList)
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, ")")
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

func (c *CompilationEngine) compileParameterList() (*Token, error) {
	paramList := &Token{
		Type:     parameterList,
		Children: make([]*Token, 0),
	}
	t, ok := c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, ")")
	}
	c.rewind()
	if t.Value == ")" && t.Type == symbol {
		// no parameter list values
		return paramList, nil
	}
	// [type] must be keyword or identity
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "type")
	}
	if !isType(t) {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, "type")
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
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "varDec | statements")
	}
	c.rewind()
	if t.Value == "var" && t.Type == keyword {
		varDecTokens, err := c.compileVarDec()
		if err != nil {
			return nil, err
		}
		subrBody.Children = append(subrBody.Children, varDecTokens)
	}

	// statements
	statements, err := c.compileStatements()
	if err != nil {
		return nil, err
	}
	subrBody.Children = append(subrBody.Children, statements)

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

// ------------- statements -------------

func (c *CompilationEngine) compileStatements() (*Token, error) {
	statements := &Token{
		Type:     statements,
		Children: make([]*Token, 0),
	}

	// ["let" | "if" | "while" | "do" | "return"] must be keyword "let" or "if" or "while" or "do" or "return"
LOOP:
	for {
		t, ok := c.next()
		if !ok {
			return nil, fmt.Errorf(ErrNoTokenExists, "let | if | while | do | return")
		}
		switch t.Value {
		case "let":
			c.rewind()
			t, err := c.compileLetStatements()
			if err != nil {
				return nil, err
			}
			statements.Children = append(statements.Children, t)
		case "if":
			c.rewind()
			t, err := c.compileIfStatements()
			if err != nil {
				return nil, err
			}
			statements.Children = append(statements.Children, t)
		case "while":
			c.rewind()
			t, err := c.compileWhileStatements()
			if err != nil {
				return nil, err
			}
			statements.Children = append(statements.Children, t)
		case "do":
			c.rewind()
			t, err := c.compileDoStatements()
			if err != nil {
				return nil, err
			}
			statements.Children = append(statements.Children, t)
		case "return":
			c.rewind()
			t, err := c.compileReturnStatements()
			if err != nil {
				return nil, err
			}
			statements.Children = append(statements.Children, t)
		default:
			c.rewind()
			break LOOP
		}
	}
	return statements, nil
}

func (c *CompilationEngine) compileLetStatements() (*Token, error) {
	// ["let"] must be keyword "let"
	t, ok := c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "let")
	}
	if t.Value != "let" || t.Type != keyword {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, "let")
	}
	letStatement := &Token{
		Type:     letStatement,
		Children: make([]*Token, 0),
	}
	letStatement.Children = append(letStatement.Children, t)

	// [varName] must be identifier
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "varName")
	}
	if t.Type != identifier {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, identifier)
	}
	letStatement.Children = append(letStatement.Children, t)

	// ["["] must be "["
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "[ | =")
	}
	if t.Value == "[" && t.Type == symbol {
		letStatement.Children = append(letStatement.Children, t)

		// compile expression
		expr, err := c.compileExpression()
		if err != nil {
			return nil, err
		}
		letStatement.Children = append(letStatement.Children, expr)

		// ["]"] must be "]"
		t, ok = c.next()
		if !ok {
			return nil, fmt.Errorf(ErrNoTokenExists, "]")
		}
		if t.Value != "]" || t.Type != symbol {
			return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
		}
		letStatement.Children = append(letStatement.Children, t)
		// ["="] must be "="
		t, ok = c.next()
		if !ok {
			return nil, fmt.Errorf(ErrNoTokenExists, "=")
		}
	}
	if t.Value != "=" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	letStatement.Children = append(letStatement.Children, t)

	// compile expression
	expr, err := c.compileExpression()
	if err != nil {
		return nil, err
	}
	letStatement.Children = append(letStatement.Children, expr)

	// [";"] must be ";"
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, ";")
	}
	if t.Value != ";" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	letStatement.Children = append(letStatement.Children, t)

	return letStatement, nil
}

func (c *CompilationEngine) compileIfStatements() (*Token, error) {
	// ["if"] must be keyword "if"
	t, ok := c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "if")
	}
	if t.Value != "if" || t.Type != keyword {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, "if")
	}
	ifStatement := &Token{
		Type:     ifStatement,
		Children: make([]*Token, 0),
	}
	ifStatement.Children = append(ifStatement.Children, t)

	// ["("] must be "("
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "(")
	}
	if t.Value != "(" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	ifStatement.Children = append(ifStatement.Children, t)

	// compile expression
	expr, err := c.compileExpression()
	if err != nil {
		return nil, err
	}
	ifStatement.Children = append(ifStatement.Children, expr)

	// [")"] must be ")"
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, ")")
	}
	if t.Value != ")" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	ifStatement.Children = append(ifStatement.Children, t)

	// ["{"] must be "{"
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "{")
	}
	if t.Value != "{" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	ifStatement.Children = append(ifStatement.Children, t)

	// statements
	statements, err := c.compileStatements()
	if err != nil {
		return nil, err
	}
	ifStatement.Children = append(ifStatement.Children, statements)

	// ["}"] must be "}"
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "}")
	}
	if t.Value != "}" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	ifStatement.Children = append(ifStatement.Children, t)

	// if exists ["else"], parse "else" section
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "else")
	}
	if t.Value != "else" && t.Type != keyword {
		// end of if statement
		c.rewind()
		return ifStatement, nil
	}
	ifStatement.Children = append(ifStatement.Children, t)

	// ["{"] must be "{"
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "{")
	}
	if t.Value != "{" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	ifStatement.Children = append(ifStatement.Children, t)

	// statements
	statements, err = c.compileStatements()
	if err != nil {
		return nil, err
	}
	ifStatement.Children = append(ifStatement.Children, statements)

	// ["}"] must be "}"
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "}")
	}
	if t.Value != "}" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	ifStatement.Children = append(ifStatement.Children, t)

	return ifStatement, nil
}

func (c *CompilationEngine) compileWhileStatements() (*Token, error) {
	// ["if"] must be keyword "if"
	t, ok := c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "if")
	}
	if t.Value != "if" || t.Type != keyword {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, "if")
	}
	whileStatement := &Token{
		Type:     whileStatement,
		Children: make([]*Token, 0),
	}
	whileStatement.Children = append(whileStatement.Children, t)

	// ["("] must be "("
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "(")
	}
	if t.Value != "(" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	whileStatement.Children = append(whileStatement.Children, t)

	// compile expression
	expr, err := c.compileExpression()
	if err != nil {
		return nil, err
	}
	whileStatement.Children = append(whileStatement.Children, expr)

	// [")"] must be ")"
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, ")")
	}
	if t.Value != ")" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	whileStatement.Children = append(whileStatement.Children, t)

	// ["{"] must be "{"
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "{")
	}
	if t.Value != "{" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	whileStatement.Children = append(whileStatement.Children, t)

	// statements
	statements, err := c.compileStatements()
	if err != nil {
		return nil, err
	}
	whileStatement.Children = append(whileStatement.Children, statements)

	// ["}"] must be "}"
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "}")
	}
	if t.Value != "}" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	whileStatement.Children = append(whileStatement.Children, t)
	return whileStatement, nil
}

func (c *CompilationEngine) compileDoStatements() (*Token, error) {
	// ["do"] must be keyword "do"
	t, ok := c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "do")
	}
	if t.Value != "do" || t.Type != keyword {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, "do")
	}
	doStatement := &Token{
		Type:     doStatement,
		Children: make([]*Token, 0),
	}
	doStatement.Children = append(doStatement.Children, t)

	// subroutine call
	if err := c.compileSubroutineCall(doStatement); err != nil {
		return nil, err
	}

	// [";"] must be ";"
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, ";")
	}
	if t.Value != ";" || t.Type != symbol {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
	}
	doStatement.Children = append(doStatement.Children, t)

	return doStatement, nil
}

func (c *CompilationEngine) compileReturnStatements() (*Token, error) {
	// ["return"] must be keyword "return"
	t, ok := c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "return")
	}
	if t.Value != "return" || t.Type != keyword {
		return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, "return")
	}
	returnStatement := &Token{
		Type:     returnStatement,
		Children: make([]*Token, 0),
	}
	returnStatement.Children = append(returnStatement.Children, t)

	// [";" | expression] must be ";" or expression
	t, ok = c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, ";")
	}
	if t.Value != ";" {
		// compile expression
		expr, err := c.compileExpression()
		if err != nil {
			return nil, err
		}
		returnStatement.Children = append(returnStatement.Children, expr)
	}
	returnStatement.Children = append(returnStatement.Children, t)

	return returnStatement, nil
}

// ------------- expression -------------

func (c *CompilationEngine) compileExpression() (*Token, error) {
	expr := &Token{
		Type:     expression,
		Children: make([]*Token, 0),
	}

	// compile term
	term, err := c.compileTerm()
	if err != nil {
		return nil, err
	}
	expr.Children = append(expr.Children, term)

	// if exists ["op"], parse "op" section
	t, ok := c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "op")
	}
	if !isOP(t) {
		c.rewind()
		return expr, nil
	}
	expr.Children = append(expr.Children, t)

	// compile term
	term, err = c.compileTerm()
	if err != nil {
		return nil, err
	}
	expr.Children = append(expr.Children, term)

	return expr, nil
}

func (c *CompilationEngine) compileTerm() (*Token, error) {
	term := &Token{
		Type:     term,
		Children: make([]*Token, 0),
	}

	t, ok := c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, "term")
	}
	switch t.Value {
	case "true", "false", "null", "this":
		// keywordConstant
		if t.Type != keyword {
			return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, "true | false | null | this")
		}
		term.Children = append(term.Children, t)
	case "-", "~":
		// unaryOp
		if t.Type != symbol {
			return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, "- | ~")
		}
		term.Children = append(term.Children, t)
		// compile term
		innerTerm, err := c.compileTerm()
		if err != nil {
			return nil, err
		}
		term.Children = append(term.Children, innerTerm)
	case "[":
		if t.Type != symbol {
			return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
		}
		term.Children = append(term.Children, t)
		// compile expression
		expr, err := c.compileExpression()
		if err != nil {
			return nil, err
		}
		term.Children = append(term.Children, expr)
		// ["]"] must be "]"
		t, ok = c.next()
		if !ok {
			return nil, fmt.Errorf(ErrNoTokenExists, "]")
		}
		if t.Value != "]" || t.Type != symbol {
			return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
		}
		term.Children = append(term.Children, t)
	case "(":
		if t.Type != symbol {
			return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
		}
		term.Children = append(term.Children, t)
		// compile expression
		expr, err := c.compileExpression()
		if err != nil {
			return nil, err
		}
		term.Children = append(term.Children, expr)
		// [")"] must be ")"
		t, ok = c.next()
		if !ok {
			return nil, fmt.Errorf(ErrNoTokenExists, ")")
		}
		if t.Value != ")" || t.Type != symbol {
			return nil, fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
		}
		term.Children = append(term.Children, t)
	default:
		if isIntegerConstant(t) || isStringConstant(t) {
			term.Children = append(term.Children, t)
		} else if t.Type == identifier {
			// subroutine call
			c.rewind()
			if err := c.compileSubroutineCall(term); err != nil {
				return nil, err
			}
		}
	}

	return term, nil
}

// FIXME subroutineName(expressionList) | (className | varName).subroutineName(expressionList)
func (c *CompilationEngine) compileSubroutineCall(parent *Token) error {
	// [subroutineName | (className | varName)] must be identifier
	t, ok := c.next()
	if !ok {
		return fmt.Errorf(ErrNoTokenExists, "subroutineName | (className | varName)")
	}
	if t.Type != identifier {
		return fmt.Errorf(ErrNotSuitableToken, t.Value, identifier)
	}
	parent.Children = append(parent.Children, t)

	// ["(" | "."] must be "(" or "."
	t, ok = c.next()
	if !ok {
		return fmt.Errorf(ErrNoTokenExists, "( | .")
	}

	if t.Value == "(" && t.Type == symbol {
		// ["("] must be "("
		parent.Children = append(parent.Children, t)
		// expresssion list
		exprList, err := c.compileExpressionList()
		if err != nil {
			return err
		}
		parent.Children = append(parent.Children, exprList)

		// [")"] must be ")"
		t, ok = c.next()
		if !ok {
			return fmt.Errorf(ErrNoTokenExists, ")")
		}
		if t.Value != ")" || t.Type != symbol {
			return fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
		}
		parent.Children = append(parent.Children, t)
	} else if t.Value == "." && t.Type == symbol {
		// ["."] must be "."
		parent.Children = append(parent.Children, t)
		// [subroutineName] must be identifier
		t, ok = c.next()
		if !ok {
			return fmt.Errorf(ErrNoTokenExists, "subroutineName")
		}
		if t.Type != identifier {
			return fmt.Errorf(ErrNotSuitableToken, t.Value, identifier)
		}
		parent.Children = append(parent.Children, t)

		// ["("] must be "("
		t, ok = c.next()
		if !ok {
			return fmt.Errorf(ErrNoTokenExists, "(")
		}
		if t.Value != "(" || t.Type != symbol {
			return fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
		}
		parent.Children = append(parent.Children, t)

		// expresssion list
		exprList, err := c.compileExpressionList()
		if err != nil {
			return err
		}
		parent.Children = append(parent.Children, exprList)

		// [")"] must be ")"
		t, ok = c.next()
		if !ok {
			return fmt.Errorf(ErrNoTokenExists, ")")
		}
		if t.Value != ")" || t.Type != symbol {
			return fmt.Errorf(ErrNotSuitableToken, t.Value, symbol)
		}
		parent.Children = append(parent.Children, t)

	} else {
		return fmt.Errorf(ErrNotSuitableToken, t.Value, "( | .)")
	}

	return nil
}

func (c *CompilationEngine) compileExpressionList() (*Token, error) {
	exprList := &Token{
		Type:     expressionList,
		Children: make([]*Token, 0),
	}
	t, ok := c.next()
	if !ok {
		return nil, fmt.Errorf(ErrNoTokenExists, ")")
	}
	c.rewind()
	if t.Value == ")" && t.Type == symbol {
		// no expression list values
		return exprList, nil
	}
	// compile expression
	expr, err := c.compileExpression()
	if err != nil {
		return nil, err
	}
	exprList.Children = append(exprList.Children, expr)
	// if exists [","], parse "," next expression list
	for {
		t, ok := c.next()
		if !ok {
			return nil, fmt.Errorf(ErrNoTokenExists, ",")
		}
		if t.Value != "," {
			c.rewind()
			return exprList, nil
		}
		exprList.Children = append(exprList.Children, t)

		// compile expression
		expr, err := c.compileExpression()
		if err != nil {
			return nil, err
		}
		exprList.Children = append(exprList.Children, expr)
	}
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

func isStringConstant(t *Token) bool {
	if t.Type != stringConstant {
		return false
	}
	for _, r := range t.Value {
		if string(r) == "\"" || unicode.IsControl(r) {
			return false
		}
	}
	return true
}

func isIntegerConstant(t *Token) bool {
	if t.Type != integerConstant {
		return false
	}
	for _, r := range t.Value {
		if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

var op = []string{
	"+",
	"-",
	"*",
	"/",
	"&",
	"|",
	"<",
	">",
	"=",
}

func isOP(t *Token) bool {
	for _, o := range op {
		if t.Value == o && t.Type == symbol {
			return true
		}
	}
	return false
}
