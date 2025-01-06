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

	// cwriter.InitSP()
	cwriter.Comment(fmt.Sprintf("---%s---", source.Name()))

	for scanner.Scan() {
		line := scanner.Text()
		commandLine := strings.TrimLeft(line, " ")
		if isComment(commandLine) || commandLine == "" {
			continue
		}
		command := strings.ToLower(strings.Split(commandLine, " ")[0])

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
			seg, index := parsePushPop(line)
			cwriter.Push(seg, index)
		case popCommand:
			seg, index := parsePushPop(line)
			cwriter.Pop(seg, index)
		case labelCommand:
		case gotoCommand:
		case ifCommand:
		case functionCommand:
		case returnCommand:
		case callCommand:
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
	"if":       ifCommand,
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

func isComment(l string) bool {
	return strings.HasPrefix(l, "//")
}

func parsePushPop(l string) (seg Segment, index int) {
	commands := strings.Split(l, " ")
	if len(commands) != 3 {
		log.Fatalf("invalid push command %q has detected\n", l)
	}
	index, err := strconv.Atoi(commands[2])
	if err != nil {
		log.Fatalf("invalid constant value %q has detected: %v\n", commands[2], err)
	}
	seg, ok := Segments[commands[1]]
	if !ok {
		log.Fatalf("invalid segment value %q has detected\n", commands[1])
	}
	return seg, index
}

func (p Parser) Close() error {
	if err := p.dest.Close(); err != nil {
		log.Fatalln(err.Error())
	}
	return nil
}
