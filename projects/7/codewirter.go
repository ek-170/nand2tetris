package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"
)

const max15BitInt = 32767
const initFuncName = "Sys.init"

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

var symbols = map[string]int{
	"END_EQ": 1,
	"END_LT": 1,
	"END_GT": 1,
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
}

func NewCodeWriter(srcFileName string, wr io.WriteCloser) CodeWriter {
	writer := bufio.NewWriterSize(wr, 1048576) // default is 1MiB

	return CodeWriter{
		wr:          writer,
		srcFileName: srcFileName,
		newLineChar: "\n",
		sb:          strings.Builder{},
	}
}

func (cw *CodeWriter) InitSP() {
	// SP = 256
	cw.writeLine(fmt.Sprintf("@%v", minSP))
	cw.writeLine("D=A")
	cw.writeLine("@SP")
	cw.writeLine("M=D")
	cw.Call(initFuncName, 0)
}

func (cw *CodeWriter) Comment(c string) {
	cw.writeLine("// " + c)
}

// --- Arithmetic ---

func (cw *CodeWriter) Add() {
	cw.Comment("add")

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
	cw.Comment("sub")

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
	cw.Comment("eq")

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

	count, ok := symbols["END_EQ"]
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
	symbols["END_EQ"] = count + 1
}

func (cw *CodeWriter) Lt() {
	cw.Comment("lt")

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

	count, ok := symbols["END_LT"]
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
	symbols["END_LT"] = count + 1
}

func (cw *CodeWriter) Gt() {
	cw.Comment("gt")

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

	count, ok := symbols["END_GT"]
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
	symbols["END_GT"] = count + 1
}

func (cw *CodeWriter) And() {
	cw.Comment("and")
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
	cw.Comment("or")
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
	cw.Comment("neg")
	cw.writeLine("@SP")
	cw.writeLine("A=M-1")
	cw.writeLine("M=-M")
}

func (cw *CodeWriter) Not() {
	cw.Comment("not")
	cw.writeLine("@SP")
	cw.writeLine("A=M-1")
	cw.writeLine("M=!M")
}

// --- Push / Pop ---

func (cw *CodeWriter) Push(seg Segment, index int) {
	if index > max15BitInt || index < 0 {
		log.Fatalf("invalid constant value %q has detected, max is %v, min is %v\n", index, max15BitInt, 0)
	}
	cw.Comment(fmt.Sprintf("push %v %v", seg, index))

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
	cw.Comment(fmt.Sprintf("pop %v %v", seg, index))

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

// --- Flow ---

func (cw *CodeWriter) Label(l string) {
	cw.Comment(fmt.Sprintf("label %v", l))

	if len(l) == 0 {
		log.Fatalf("invalid label %q", l)
	}
	first := rune(l[0])
	if !(unicode.IsLetter(first) || first == '_' || first == '.' || first == ':') {
		log.Fatalf("invalid label %q", l)
	}
	for _, r := range l[1:] {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' || r == ':') {
			log.Fatalf("invalid label %q", l)
		}
	}
	cw.writeLine("(" + l + ")")
}

func (cw *CodeWriter) Goto(l string) {
	cw.Comment(fmt.Sprintf("goto %v", l))

	cw.writeLine("@" + l)
	cw.writeLine("0;JMP")
}

func (cw *CodeWriter) IfGoto(l string) {
	cw.Comment(fmt.Sprintf("if-goto %v", l))

	cw.writeLine("@SP")
	cw.writeLine("AM=M-1") // decrement SP
	cw.writeLine("D=M")    // value in SP
	cw.writeLine("@" + l)
	cw.writeLine("D;JNE")
}

// --- Function ---

func (cw *CodeWriter) Func(name string, local int) {
	cw.Comment(fmt.Sprintf("function %v %v", name, local))

	cw.Label(name)
	if local == 0 {
		return
	}
	cw.writeLine("@LCL")
	cw.writeLine("AD=M")
	if local == 1 {
		cw.writeLine("M=0")
	} else if local > 1 {
		cw.writeLine("M=0")
		for i := 0; i < local-1; i++ {
			cw.writeLine("AD=D+1")
			cw.writeLine("M=0")
		}
	}
	cw.writeLine("@SP")
	cw.writeLine("M=D+1")
}

func (cw *CodeWriter) Call(name string, arg int) {
	cw.Comment(fmt.Sprintf("call %v %v", name, arg))

	// push return-address
	var returnLabel string
	s, ok := symbols[name+".return"]
	if ok {
		returnLabel = fmt.Sprintf("%v.return%v", name, s+1)
		symbols[name+".return"] = s + 1
	} else {
		returnLabel = fmt.Sprintf("%v.return%v", name, 1)
		symbols[name+".return"] = 1
	}
	cw.writeLine("@" + returnLabel)
	cw.writeLine("D=A")
	cw.writeLine("@SP")
	cw.writeLine("AM=M+1") // increment SP
	cw.writeLine("A=A-1")
	cw.writeLine("M=D") // push value
	// push current LCL
	cw.writeLine("@LCL")
	cw.writeLine("D=M")
	cw.writeLine("@SP")
	cw.writeLine("AM=M+1") // increment SP
	cw.writeLine("A=A-1")
	cw.writeLine("M=D") // push value
	// push current ARG
	cw.writeLine("@ARG")
	cw.writeLine("D=M")
	cw.writeLine("@SP")
	cw.writeLine("AM=M+1") // increment SP
	cw.writeLine("A=A-1")
	cw.writeLine("M=D") // push value
	// push current THIS
	cw.writeLine("@THIS")
	cw.writeLine("D=M")
	cw.writeLine("@SP")
	cw.writeLine("AM=M+1") // increment SP
	cw.writeLine("A=A-1")
	cw.writeLine("M=D") // push value
	// push current THAT
	cw.writeLine("@THAT")
	cw.writeLine("D=M")
	cw.writeLine("@SP")
	cw.writeLine("AM=M+1") // increment SP
	cw.writeLine("A=A-1")
	cw.writeLine("M=D") // push value
	// ARG = SP-n-5
	cw.writeLine("@5")
	cw.writeLine("D=A")
	cw.writeLine(fmt.Sprintf("@%v", arg))
	cw.writeLine("D=D+A")
	cw.writeLine("@SP")
	cw.writeLine("D=M-D")
	cw.writeLine("@ARG")
	cw.writeLine("M=D")
	// LCL = SP
	cw.writeLine("@SP")
	cw.writeLine("D=M")
	cw.writeLine("@LCL")
	cw.writeLine("M=D")
	// goto f
	cw.writeLine("@" + name)
	cw.writeLine("0;JMP")
	// (return-address)
	cw.Label(returnLabel)
}

func (cw *CodeWriter) Return() {
	cw.Comment("return")

	// FRAME = LCL
	cw.writeLine("@LCL")
	cw.writeLine("D=M")
	cw.writeLine("@R13")
	cw.writeLine("M=D")
	// RET(R14) = *(FRAME - 5)
	cw.writeLine("@5")
	cw.writeLine("D=A")
	cw.writeLine("@R13")
	cw.writeLine("A=M-D")
	cw.writeLine("D=M")
	cw.writeLine("@R14")
	cw.writeLine("M=D")
	// *ARG = Result
	cw.writeLine("@SP")
	cw.writeLine("A=M-1")
	cw.writeLine("D=M")
	cw.writeLine("@ARG")
	cw.writeLine("A=M")
	cw.writeLine("M=D")
	// SP = ARG + 1
	cw.writeLine("@ARG")
	cw.writeLine("D=M+1")
	cw.writeLine("@SP")
	cw.writeLine("M=D")
	// THAT = *(FRAME-1)
	cw.writeLine("@R13")
	cw.writeLine("A=M-1")
	cw.writeLine("D=M")
	cw.writeLine("@THAT")
	cw.writeLine("M=D")
	// THIS = *(FRAME-2)
	cw.writeLine("@2")
	cw.writeLine("D=A")
	cw.writeLine("@R13")
	cw.writeLine("A=M-D")
	cw.writeLine("D=M")
	cw.writeLine("@THIS")
	cw.writeLine("M=D")
	// ARG = *(FRAME-3)
	cw.writeLine("@3")
	cw.writeLine("D=A")
	cw.writeLine("@R13")
	cw.writeLine("A=M-D")
	cw.writeLine("D=M")
	cw.writeLine("@ARG")
	cw.writeLine("M=D")
	// LCL = *(FRAME-4)
	cw.writeLine("@4")
	cw.writeLine("D=A")
	cw.writeLine("@R13")
	cw.writeLine("A=M-D")
	cw.writeLine("D=M")
	cw.writeLine("@LCL")
	cw.writeLine("M=D")
	// goto RET
	cw.writeLine("@R14")
	cw.writeLine("A=M")
	cw.writeLine("0;JMP")
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
