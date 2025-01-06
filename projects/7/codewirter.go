package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

const max15BitInt = 32767

type Segment string

const (
	local    Segment = "local"
	argument Segment = "argument"
	this     Segment = "this"
	that     Segment = "that"
	pointer  Segment = "pointer"
	temp     Segment = "temp"
	constant Segment = "constant"
	static   Segment = "static"
)

var Segments = map[string]Segment{
	"local":    local,
	"argument": argument,
	"this":     this,
	"that":     that,
	"pointer":  pointer,
	"temp":     temp,
	"constant": constant,
	"static":   static,
}

const (
	minSP    = 256
	tempBase = 5
)

type CodeWriter struct {
	wr          *bufio.Writer
	srcFileName string
	newLineChar string
	sb          strings.Builder
	symbols     map[string]int
}

func NewCodeWriter(srcFileName string, wr io.WriteCloser) CodeWriter {
	writer := bufio.NewWriterSize(wr, 1048576) // default is 1MiB

	return CodeWriter{
		wr:          writer,
		srcFileName: srcFileName,
		newLineChar: "\n",
		sb:          strings.Builder{},
		symbols: map[string]int{
			"END_EQ": 1,
			"END_LT": 1,
			"END_GT": 1,
		},
	}
}

func (cw *CodeWriter) InitSP() {
	cw.writeLine(fmt.Sprintf("@%v", minSP))
	cw.writeLine("D=A")
	cw.writeLine("@SP")
	cw.writeLine("M=D")
}

func (cw *CodeWriter) Comment(c string) {
	cw.writeLine("// " + c)
}

func (cw *CodeWriter) Add() {
	cw.writeLine("@SP")
	// decrement SP and set A Register
	cw.writeLine("AM=M-1")
	// set y to D Register
	cw.writeLine("D=M")
	// set x address to A Register
	cw.writeLine("A=A-1")
	// x = x + y
	cw.writeLine("M=M+D")
}

func (cw *CodeWriter) Sub() {
	cw.writeLine("@SP")
	// decrement SP and set A Register
	cw.writeLine("AM=M-1")
	// set y to D Register
	cw.writeLine("D=M")
	// set x address to A Register
	cw.writeLine("A=A-1")
	// x = x - y
	cw.writeLine("M=M-D")
}

func (cw *CodeWriter) Eq() {
	cw.writeLine("@SP")
	// decrement SP and set A Register
	cw.writeLine("AM=M-1")
	// set y to D Register
	cw.writeLine("D=M")
	// set x address to A Register
	cw.writeLine("A=A-1")
	cw.writeLine("D=M-D")
	// default is -1(true)
	cw.writeLine("M=-1")

	count, ok := cw.symbols["END_EQ"]
	if !ok {
		log.Fatalf("symbol %q has not found\n", "END_EQ")
	}
	// if x - y != 0, set 0(false)
	cw.writeLine(fmt.Sprintf("@END_EQ%v", count))
	cw.writeLine("D;JEQ")
	cw.writeLine("@SP")
	cw.writeLine("A=M-1")
	cw.writeLine("M=0")

	cw.writeLine(fmt.Sprintf("(END_EQ%v)", count))
	cw.symbols["END_EQ"] = count + 1
}

func (cw *CodeWriter) Lt() {
	cw.writeLine("@SP")
	// decrement SP and set A Register
	cw.writeLine("AM=M-1")
	// set y to D Register
	cw.writeLine("D=M")
	// set x address to A Register
	cw.writeLine("A=A-1")
	cw.writeLine("D=M-D")
	// default is -1(true)
	cw.writeLine("M=-1")

	count, ok := cw.symbols["END_LT"]
	if !ok {
		log.Fatalf("symbol %q has not found\n", "END_LT")
	}
	// if NOT x - y < 0, set 0(false)
	cw.writeLine(fmt.Sprintf("@END_LT%v", count))
	cw.writeLine("D;JLT")
	cw.writeLine("@SP")
	cw.writeLine("A=M-1")
	cw.writeLine("M=0")

	cw.writeLine(fmt.Sprintf("(END_LT%v)", count))
	cw.symbols["END_LT"] = count + 1
}

func (cw *CodeWriter) Gt() {
	cw.writeLine("@SP")
	// decrement SP and set A Register
	cw.writeLine("AM=M-1")
	// set y to D Register
	cw.writeLine("D=M")
	// set x address to A Register
	cw.writeLine("A=A-1")
	cw.writeLine("D=M-D")
	// default is -1(true)
	cw.writeLine("M=-1")

	count, ok := cw.symbols["END_GT"]
	if !ok {
		log.Fatalf("symbol %q has not found\n", "END_GT")
	}
	// if NOT x - y > 0, set 0(false)
	cw.writeLine(fmt.Sprintf("@END_GT%v", count))
	cw.writeLine("D;JGT")
	cw.writeLine("@SP")
	cw.writeLine("A=M-1")
	cw.writeLine("M=0")

	cw.writeLine(fmt.Sprintf("(END_GT%v)", count))
	cw.symbols["END_GT"] = count + 1
}

func (cw *CodeWriter) And() {
	cw.writeLine("@SP")
	// decrement SP and set A Register
	cw.writeLine("AM=M-1")
	// set y to D Register
	cw.writeLine("D=M")
	// set x address to A Register
	cw.writeLine("A=A-1")
	cw.writeLine("M=M&D")
}

func (cw *CodeWriter) Or() {
	cw.writeLine("@SP")
	// decrement SP and set A Register
	cw.writeLine("AM=M-1")
	// set y to D Register
	cw.writeLine("D=M")
	// set x address to A Register
	cw.writeLine("A=A-1")
	cw.writeLine("M=M|D")
}

func (cw *CodeWriter) Neg() {
	cw.writeLine("@SP")
	cw.writeLine("A=M-1")
	cw.writeLine("M=-M")
}

func (cw *CodeWriter) Not() {
	cw.writeLine("@SP")
	cw.writeLine("A=M-1")
	cw.writeLine("M=!M")
}

func (cw *CodeWriter) Push(seg Segment, index int) {
	if index > max15BitInt || index < 0 {
		log.Fatalf("invalid constant value %q has detected, max is %v, min is %v\n", index, max15BitInt, 0)
	}
	switch seg {
	case local:
		cw.writeLine(fmt.Sprintf("@%v", index))
		cw.writeLine("D=A")
		cw.writeLine("@LCL")
		cw.writeLine("A=M+D") // base + index
		cw.writeLine("D=M")
	case argument:
		cw.writeLine(fmt.Sprintf("@%v", index))
		cw.writeLine("D=A")
		cw.writeLine("@ARG")
		cw.writeLine("A=M+D") // base + index
		cw.writeLine("D=M")
	case this:
		cw.writeLine(fmt.Sprintf("@%v", index))
		cw.writeLine("D=A")
		cw.writeLine("@THIS")
		cw.writeLine("A=M+D") // base + index
		cw.writeLine("D=M")
	case that:
		cw.writeLine(fmt.Sprintf("@%v", index))
		cw.writeLine("D=A")
		cw.writeLine("@THAT")
		cw.writeLine("A=M+D") // base + index
		cw.writeLine("D=M")
	case pointer:
		if index == 0 {
			cw.writeLine("@THIS")
		} else if index == 1 {
			cw.writeLine("@THAT")
		} else {
			log.Fatalf("temp index must be 0, 1, but %v", index)
		}
		cw.writeLine("D=M")
	case temp:
		if index > 7 {
			log.Fatalf("temp index must be 0 ~ 7, but %v", index)
		}
		cw.writeLine(fmt.Sprintf("@R%v", tempBase+index))
		cw.writeLine("D=M")
	case constant:
		cw.writeLine(fmt.Sprintf("@%v", index))
		cw.writeLine("D=A")
	case static:
		cw.writeLine(fmt.Sprintf("@%v.%v", cw.srcFileName, index))
		cw.writeLine("D=M")
	default:
		log.Fatalf("unknown memory segment %q has detected\n", seg)
	}
	cw.writeLine("@SP")
	cw.writeLine("AM=M+1") // increment SP
	cw.writeLine("A=A-1")
	cw.writeLine("M=D") // push value
}

func (cw *CodeWriter) Pop(seg Segment, index int) {
	if index > max15BitInt || index < 0 {
		log.Fatalf("invalid constant value %q has detected, max is %v, min is %v\n", index, max15BitInt, 0)
	}

	switch seg {
	case local:
		cw.writeLine(fmt.Sprintf("@%v", index))
		cw.writeLine("D=A")
		cw.writeLine("@LCL")
		cw.writeLine("D=M+D") // D is base + index

		cw.writeLine("@R13")
		cw.writeLine("M=D")

		cw.writeLine("@SP")
		cw.writeLine("AM=M-1") // decrement SP
		cw.writeLine("D=M")    // value in SP
		cw.writeLine("@R13")
		cw.writeLine("A=M")
		cw.writeLine("M=D")
	case argument:
		cw.writeLine(fmt.Sprintf("@%v", index))
		cw.writeLine("D=A")
		cw.writeLine("@ARG")
		cw.writeLine("D=M+D") // D is base + index

		cw.writeLine("@R13")
		cw.writeLine("M=D")

		cw.writeLine("@SP")
		cw.writeLine("AM=M-1") // decrement SP
		cw.writeLine("D=M")    // value in SP
		cw.writeLine("@R13")
		cw.writeLine("A=M")
		cw.writeLine("M=D")
	case this:
		cw.writeLine(fmt.Sprintf("@%v", index))
		cw.writeLine("D=A")
		cw.writeLine("@THIS")
		cw.writeLine("D=M+D") // D is base + index

		cw.writeLine("@R13")
		cw.writeLine("M=D")

		cw.writeLine("@SP")
		cw.writeLine("AM=M-1") // decrement SP
		cw.writeLine("D=M")    // value in SP
		cw.writeLine("@R13")
		cw.writeLine("A=M")
		cw.writeLine("M=D")
	case that:
		cw.writeLine(fmt.Sprintf("@%v", index))
		cw.writeLine("D=A")
		cw.writeLine("@THAT")
		cw.writeLine("D=M+D") // D is base + index

		cw.writeLine("@R13")
		cw.writeLine("M=D")

		cw.writeLine("@SP")
		cw.writeLine("AM=M-1") // decrement SP
		cw.writeLine("D=M")    // value in SP
		cw.writeLine("@R13")
		cw.writeLine("A=M")
		cw.writeLine("M=D")
	case pointer:
		if index != 0 && index != 1 {
			log.Fatalf("pointer index must be 0 or 1, but %v", index)
		}
		cw.writeLine("@SP")
		cw.writeLine("AM=M-1") // decrement SP
		cw.writeLine("D=M")
		if index == 0 {
			cw.writeLine("@THIS")
		} else if index == 1 {
			cw.writeLine("@THAT")
		}
		cw.writeLine("M=D")
	case temp:
		if index > 7 {
			log.Fatalf("temp index must be 0 ~ 7, but %v", index)
		}
		cw.writeLine("@SP")
		cw.writeLine("AM=M-1") // decrement SP
		cw.writeLine("D=M")    // value in SP
		cw.writeLine(fmt.Sprintf("@R%v", tempBase+index))
		cw.writeLine("M=D")
	// case constant:
	case static:
		cw.writeLine("@SP")
		cw.writeLine("AM=M-1") // decrement SP
		cw.writeLine("D=M")    // value in SP
		cw.writeLine(fmt.Sprintf("@%v.%v", cw.srcFileName, index))
		cw.writeLine("M=D")
	default:
		log.Fatalf("unknown memory segment %q has detected\n", seg)
	}
}

func (cw *CodeWriter) write(s string) {
	if _, err := cw.sb.WriteString(s); err != nil {
		log.Fatalf("failed to write to string builder: %v\n", err)
	}
}

func (cw *CodeWriter) writeLine(s string) {
	if _, err := cw.sb.WriteString(s); err != nil {
		log.Fatalf("failed to write to string builder: %v\n", err)
	}
	if _, err := cw.sb.WriteString(cw.newLineChar); err != nil {
		log.Fatalf("failed to write to string builder: %v\n", err)
	}
}

func (cw *CodeWriter) Flush() {
	if _, err := cw.wr.WriteString(cw.sb.String()); err != nil {
		log.Fatalf("failed to flush: %v\n", err)
	}
	if err := cw.wr.Flush(); err != nil {
		log.Fatalf("failed to flush: %v\n", err)
	}
}
