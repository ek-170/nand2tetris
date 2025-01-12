package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type Parser struct {
	dest io.WriteCloser
}

func NewParser(source *os.File, dest io.WriteCloser) Parser {
	return Parser{
		dest: dest,
	}
}

func NewParserWithFile(destPath string) Parser {
	dest, err := os.Create(destPath)
	if err != nil {
		log.Fatalln("could not create new file")
	}
	return Parser{
		dest: dest,
	}
}

func (p Parser) Do(source *os.File) {
	scanner := bufio.NewScanner(source)
	cwriter := NewCodeWriter(strings.TrimRight(filepath.Base(source.Name()), ".vm"), p.dest)

	cwriter.InitSP()
	cwriter.Comment(fmt.Sprintf("---%s---", source.Name()))

	for scanner.Scan() {
		line := scanner.Text()
		tokens, skip := tokenizeCommand(line)
		if skip {
			continue
		}
		if len(tokens) == 0 {
			log.Fatalf("invalid command line detected %q", tokens)
		}
		command := strings.ToLower(tokens[0])
		switch p.CommandType(command) {
		case arithCommand:
			switch command {
			case "add":
				cwriter.Add()
			case "sub":
				cwriter.Sub()
			case "neg":
				cwriter.Neg()
			case "eq":
				cwriter.Eq()
			case "gt":
				cwriter.Gt()
			case "lt":
				cwriter.Lt()
			case "and":
				cwriter.And()
			case "or":
				cwriter.Or()
			case "not":
				cwriter.Not()
			default:
				log.Fatalf("unknown arithmetic command %q\n", command)
			}
		case pushCommand:
			cwriter.Push(validatePushAndPop(tokens))
		case popCommand:
			cwriter.Pop(validatePushAndPop(tokens))
		case labelCommand:
			cwriter.Label(validateLabel(tokens))
		case gotoCommand:
			cwriter.Goto(validateProgramFlow(tokens))
		case ifCommand:
			cwriter.IfGoto(validateProgramFlow(tokens))
		case functionCommand:
			cwriter.Func(validateFuncAndCall(tokens))
		case returnCommand:
			cwriter.Return()
		case callCommand:
			cwriter.Call(validateFuncAndCall(tokens))
		default:
			log.Fatalf("unknown command %q\n", command)
		}
	}
	cwriter.Flush()
}

type cType string

const (
	arithCommand    cType = "C_ARITHMETIC"
	pushCommand     cType = "C_PUSH"
	popCommand      cType = "C_POP"
	labelCommand    cType = "C_LABEL"
	gotoCommand     cType = "C_GOTO"
	ifCommand       cType = "C_IF"
	functionCommand cType = "C_FUNCTION"
	returnCommand   cType = "C_RETURN"
	callCommand     cType = "C_CALL"
)

var commands map[string]cType = map[string]cType{
	"add":      arithCommand,
	"sub":      arithCommand,
	"neg":      arithCommand,
	"eq":       arithCommand,
	"gt":       arithCommand,
	"lt":       arithCommand,
	"and":      arithCommand,
	"or":       arithCommand,
	"not":      arithCommand,
	"push":     pushCommand,
	"pop":      popCommand,
	"label":    labelCommand,
	"goto":     gotoCommand,
	"if-goto":  ifCommand,
	"function": functionCommand,
	"return":   returnCommand,
	"call":     callCommand,
}

func (p Parser) CommandType(c string) cType {
	command, ok := commands[c]
	if !ok {
		log.Fatalf("unknown command %q has detected\n", c)
	}
	return command
}

func tokenizeCommand(line string) (tokens []string, skip bool) {
	// 行頭・行末の空白をトリム
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "//") {
		return nil, true // 空行またはコメント行はスキップ
	}

	var tokenBuilder strings.Builder
	var tokensCollected []string
	var isFirstToken, isPrevSlash bool
	isFirstToken = true

	for i, r := range line {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ':' || r == '.' || r == '_' || (isFirstToken && r == '-') {
			if isPrevSlash {
				log.Fatalf("invalid character %q detected in line: %q", r, line[:i+1])
			}
			tokenBuilder.WriteRune(r) // 有効な文字をトークンに追加
			isPrevSlash = false
			continue
		}

		if r == '/' {
			if isPrevSlash { // "//" コメント部分を検出
				isPrevSlash = false
				break
			}
			isPrevSlash = true
			continue
		}

		if unicode.IsSpace(r) { // トークンの終了を検出
			if tokenBuilder.Len() > 0 {
				if isPrevSlash {
					log.Fatalf("invalid character %q detected in line: %q", r, line[:i+1])
				}
				tokensCollected = append(tokensCollected, tokenBuilder.String())
				tokenBuilder.Reset()
				isFirstToken = false
				isPrevSlash = false
			}
			continue
		}

		// 不正な文字を検出
		log.Fatalf("invalid character %q detected in line: %q", r, line[:i+1])
	}

	if isPrevSlash {
		log.Fatalf("invalid character %q detected in line: %q", "/", line)
	}

	// 最後のトークンを追加
	if tokenBuilder.Len() > 0 {
		tokensCollected = append(tokensCollected, tokenBuilder.String())
	}
	return tokensCollected, false
}

func validatePushAndPop(tokens []string) (seg Segment, index int) {
	if len(tokens) != 3 {
		log.Fatalf("invalid push/pop command %q has detected\n", tokens)
	}
	index, err := strconv.Atoi(tokens[2])
	if err != nil {
		log.Fatalf("invalid constant value %q has detected: %v\n", tokens[2], err)
	}
	seg, ok := Segments[tokens[1]]
	if !ok {
		log.Fatalf("invalid segment value %q has detected\n", tokens[1])
	}
	return seg, index
}

func validateLabel(tokens []string) string {
	if len(tokens) != 2 {
		log.Fatalf("invalid label command %q has detected\n", tokens)
	}
	return tokens[1]
}

func validateProgramFlow(tokens []string) (label string) {
	if len(tokens) != 2 {
		log.Fatalf("invalid program flow command %q has detected\n", tokens)
	}
	return tokens[1]
}

func validateFuncAndCall(tokens []string) (name string, local int) {
	if len(tokens) != 3 {
		log.Fatalf("invalid func/call command %q has detected\n", tokens)
	}
	local, err := strconv.Atoi(tokens[2])
	if err != nil {
		log.Fatalf("invalid value %q has detected: %v\n", tokens[2], err)
	}
	return tokens[1], local
}

func (p Parser) Close() error {
	if err := p.dest.Close(); err != nil {
		log.Fatalln(err.Error())
	}
	return nil
}
