package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Parser struct {
	source *os.File
	dest   io.Writer
}

func NewParser(source *os.File, dest io.Writer) Parser {
	return Parser{
		source: source,
		dest:   dest,
	}
}

func NewParserWithFile(source *os.File, destPath string) Parser {
	dest, err := os.Create(destPath)
	if err != nil {
		log.Fatal("could not create new file")
	}
	return Parser{
		source: source,
		dest:   dest,
	}
}

func (p Parser) Do() {
	scanner := bufio.NewScanner(p.source)
	defer func() {
		if err := p.source.Close(); err != nil {
			log.Fatalln(err.Error())
		}
		closer, ok := p.dest.(io.Closer)
		if ok {
			if err := closer.Close(); err != nil {
				log.Fatalln(err.Error())
			}
		}
	}()
	romAddress := 0

	for scanner.Scan() {
		line := scanner.Text()
		command := strings.TrimLeft(line, " ")
		if isComment(command) || command == "" {
			continue
		}
		if p.CommandType(command) == lCommand {
			symbol := line[1 : len(line)-1]
			symbolTable.AddROMEntry(symbol, romAddress)
			continue
		}
		romAddress++
	}

	_, err := p.source.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Printf("Failed to seek file: %v\n", err)
		return
	}
	scanner = bufio.NewScanner(p.source)
	writer := bufio.NewWriterSize(p.dest, 1048576) // default is 1MiB
	for scanner.Scan() {
		line := scanner.Text()
		command := strings.TrimLeft(line, " ")
		if isComment(command) || command == "" {
			continue
		}
		var binaryStr string
		switch p.CommandType(command) {
		case aCommand:
			binaryStr = parseA(command)
		case cCommand:
			binaryStr = parseC(command)
		default:
			continue
		}
		_, err := writer.WriteString(binaryStr + "\n")
		if err != nil {
			log.Fatalln("write string error")
		}
	}
	if err := writer.Flush(); err != nil {
		log.Fatalln("flush error")
	}
}

type cType string

const (
	aCommand cType = "A_COMMAND" // A命令
	cCommand cType = "C_COMMAND" // C命令
	lCommand cType = "L_COMMAND" // 擬似コマンド
)

// 現コマンドの種類を返す
func (p Parser) CommandType(l string) cType {
	if strings.HasPrefix(l, "@") {
		return aCommand
	}
	if strings.HasPrefix(l, "(") && strings.HasSuffix(l, ")") {
		return lCommand
	}
	return cCommand
}

const max15BitInt = 32767

func parseA(l string) (binary string) {
	s := l[1:]
	i, err := strconv.Atoi(s)
	if err != nil {
		symbolTable.AddRAMEntry(s)
		address, ok := symbolTable.GetAddress(s)
		if !ok {
			log.Fatalln("not found symbol's address")
		}
		i = address
	}

	if i > max15BitInt {
		log.Fatalf("too large integer detected while parsing A command: %v", i)
	}
	return transform10to2WithZeroPadding(i, 16)
}

func parseC(l string) (binary string) {
	var comp, dest, jump string

	before, after, found := strings.Cut(l, "=")
	if found {
		dest = destMnemonics[strings.TrimSpace(before)]
	} else {
		dest = destMnemonics["null"]
		// If sep does not appear in s, cut returns s, "", false.
		after = before
	}

	before, after, found = strings.Cut(after, ";")
	if found {
		jump = jumpMnemonics[strings.TrimSpace(after)]
	} else {
		jump = jumpMnemonics["null"]
	}

	comp = compMnemonics[strings.TrimSpace(before)]
	return "111" + comp + dest + jump
}

func isComment(l string) bool {
	return strings.HasPrefix(l, "//")
}

func transform10to2WithZeroPadding(v int, digit int) string {
	b := strconv.FormatInt(int64(v), 2)
	if len(b) < digit {
		var sb strings.Builder
		for i := 0; i < digit-len(b); i++ {
			_, err := sb.WriteString("0")
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		b = sb.String() + b
	} else {
		log.Fatalf("too large integer %v", v)
	}
	return b
}
