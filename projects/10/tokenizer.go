package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type TokenType string

const (
	keyword         TokenType = "keyword"
	symbol          TokenType = "symbol"
	integerConstant TokenType = "integerConstant"
	stringConstant  TokenType = "stringConstant"
	identifier      TokenType = "identifier"
)

var (
	keywords = []string{
		"class",
		"constructor",
		"function",
		"method",
		"field",
		"static",
		"var",
		"int",
		"char",
		"boolean",
		"void",
		"true",
		"false",
		"null",
		"this",
		"let",
		"do",
		"if",
		"else",
		"while",
		"return",
	}
	symbols = []string{
		"{",
		"}",
		"(",
		")",
		"[",
		"]",
		".",
		",",
		";",
		"+",
		"-",
		"*",
		"/",
		"&",
		"|",
		"<",
		">",
		"=",
		"~",
	}
)

type Token struct {
	Type     TokenType
	Value    string
	Children []Token
}

type JackTokenizer struct {
	src io.Reader
}

func NewJackTokenizer(src io.Reader) JackTokenizer {
	return JackTokenizer{src: src}
}

const (
	space = " "
	lf    = "\n"
	cr    = "\r"
	crlf  = "\r\n"
	tab   = "\t"
)

func (jt JackTokenizer) Tokenize() ([]Token, error) {

	var maybeCommentStartChars, maybeCommentEndChars string
	isLineCommenting, isMultiLineCommenting, isAPICommenting := false, false, false

	tokens := []Token{}
	builder := strings.Builder{}
	b := bufio.NewReader(jt.src)

redo:
	token := Token{}
	for {
		r, _, err := b.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		char := string(r)

		if maybeCommentStartChars == "/" && !(char == "/" || char == "*") {
			token = Token{}
			token.Type = symbol
			token.Value = "/"
			tokens = append(tokens, token)
			maybeCommentStartChars = ""
			token = Token{}
		}

		// タブ、スペースはtokenize中以外無視
		if char == space || char == tab {
			if token.Type == stringConstant {
				if _, err = builder.WriteRune(r); err != nil {
					return tokens, err
				}
				continue
			} else if token.Type != "" {
				token.Value = builder.String()
				tokens = append(tokens, token)
				builder.Reset()
				goto redo
			}
			continue
		}

		// 改行はトークンの区切りか行コメント解除
		if char == lf || char == cr || char == crlf {
			if token.Type == stringConstant {
				return nil, fmt.Errorf("%v%q has appeared in token type %q", builder.String(), char, token.Type)
			} else if token.Type != "" {
				token.Value = builder.String()
				tokens = append(tokens, token)
				builder.Reset()
				goto redo
			}
			if isLineCommenting {
				maybeCommentStartChars = ""
				isLineCommenting = false
			}
			continue
		}

		// /の場合は次の文字がコメントを示す文字か確認
		if char == "/" {
			if maybeCommentStartChars == "/" {
				isLineCommenting = true
				maybeCommentStartChars = ""
			} else if !(isLineCommenting || isMultiLineCommenting || isAPICommenting) {
				maybeCommentStartChars = "/"
			}
			if isMultiLineCommenting || isAPICommenting {
				if maybeCommentEndChars == "*" {
					maybeCommentStartChars = ""
					maybeCommentEndChars = ""
					isMultiLineCommenting = false
					isAPICommenting = false
				}
			}
			continue
		}

		if char == "*" {
			if isMultiLineCommenting || isAPICommenting {
				maybeCommentEndChars = "*"
				continue
			}
			if maybeCommentStartChars == "/" {
				maybeCommentStartChars = "/*"
				isMultiLineCommenting = true
				continue
			} else if maybeCommentStartChars == "/*" {
				isMultiLineCommenting = false
				isAPICommenting = true
				continue
			}
			token.Type = symbol
			token.Value = "*"
			tokens = append(tokens, token)
			goto redo
		}

		if isLineCommenting {
			continue
		}

		if isMultiLineCommenting || isAPICommenting {
			maybeCommentEndChars = ""
			continue
		}

		// "が出たら次に"が出るまでstringConstant開始
		if char == "\"" {
			if token.Type == stringConstant {
				token.Value = builder.String()
				tokens = append(tokens, token)
				builder.Reset()
				goto redo
			} else if token.Type == "" {
				token.Type = stringConstant
			} else {
				return nil, fmt.Errorf("%v%q has appeared in token type %q", builder.String(), char, token.Type)
			}
			continue
		}

		// 数字が出たらintegerConstant
		if unicode.IsNumber(r) {
			if token.Type == "" {
				token.Type = integerConstant
			}
			if _, err = builder.WriteRune(r); err != nil {
				return tokens, err
			}
			continue
		}

		// アルファベットで始まったらidentifierかkeyword
		if unicode.IsLetter(r) {
			if token.Type == "" {
				token.Type = identifier
			} else if token.Type == integerConstant {
				return nil, fmt.Errorf("%v%q has appeared in token type %q", builder.String(), char, token.Type)
			}
			if _, err = builder.WriteRune(r); err != nil {
				return tokens, err
			}
			isKeyword := false
			for _, k := range keywords {
				if k == builder.String() {
					token.Type = keyword
					token.Value = builder.String()
					tokens = append(tokens, token)
					builder.Reset()
					isKeyword = true
					goto redo
				}
			}
			if isKeyword {
				goto redo
			}
			continue
		}

		// symbolをハンドリング
		isSymbol := false
		for _, sym := range symbols {
			if sym == char {
				isSymbol = true
			}
		}
		if isSymbol {
			if token.Type != "" {
				token.Value = builder.String()
				tokens = append(tokens, token)
				builder.Reset()
			}
			token = Token{}
			token.Type = symbol
			token.Value = char
			tokens = append(tokens, token)
			goto redo
		}

		if token.Type == stringConstant {
			if _, err = builder.WriteRune(r); err != nil {
				return tokens, err
			}
			continue
		}

		return nil, fmt.Errorf("invalid char %q has appeared in token type %q", char, token.Type)
	}

	return tokens, nil
}
