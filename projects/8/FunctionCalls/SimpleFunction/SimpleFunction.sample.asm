// ---../8/FunctionCalls/SimpleFunction/SimpleFunction.vm---
(SimpleFunction.test)
@LCL
AD=M
M=0
AD=D+1
M=0
@SP
M=D+1
// push local 0
@0
D=A
@LCL
A=M+D
D=M
@SP
AM=M+1
A=A-1
M=D
// push local 1
@1
D=A
@LCL
A=M+D
D=M
@SP
AM=M+1
A=A-1
M=D
// add
@SP
AM=M-1
D=M
A=A-1
M=M+D
// not
@SP
A=M-1
M=!M
// push argument 0
@0
D=A
@ARG
A=M+D
D=M
@SP
AM=M+1
A=A-1
M=D
// add
@SP
AM=M-1
D=M
A=A-1
M=M+D
// push argument 1
@1
D=A
@ARG
A=M+D
D=M
@SP
AM=M+1
A=A-1
M=D
// sub
@SP
AM=M-1
D=M
A=A-1
M=M-D
//
// FRAME = LCL
@LCL
D=M
@R13
M=D // FRAME
// RET = *(FRAME - 5)
@5
D=A
@R13
A=M-D // FRAME - 5
D=M // *(FRAME - 5)
@R14 // RET
M=D
// *ARG = Result
@SP
A=M-1 // Result
D=M
@ARG
A=M
M=D
// SP = ARG + 1
@ARG
D=M+1
@SP
M=D
// THAT = *(FRAME-1)
@R13
A=M-1 // FRAME - 1
D=M // *(FRAME - 1)
@THAT
M=D
// THIS = *(FRAME-2)
@2
D=A
@R13
A=M-D // FRAME - 2
D=M // *(FRAME - 2)
@THIS
M=D
// ARG = *(FRAME-3)
@3
D=A
@R13
A=M-D // FRAME - 3
D=M // *(FRAME - 3)
@ARG
M=D
// LCL = *(FRAME-4)
@4
D=A
@R13
A=M-D // FRAME - 4
D=M // *(FRAME - 4)
@LCL
M=D
// goto RET
@R14
A=M
0;JMP
