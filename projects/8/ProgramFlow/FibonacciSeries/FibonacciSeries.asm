// ---../8/ProgramFlow/FibonacciSeries/FibonacciSeries.vm---
@1
D=A
@ARG
A=M+D
D=M
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
@THAT
M=D
@0
D=A
@SP
AM=M+1
A=A-1
M=D
@0
D=A
@THAT
D=M+D
@R13
M=D
@SP
AM=M-1
D=M
@R13
A=M
M=D
@1
D=A
@SP
AM=M+1
A=A-1
M=D
@1
D=A
@THAT
D=M+D
@R13
M=D
@SP
AM=M-1
D=M
@R13
A=M
M=D
@0
D=A
@ARG
A=M+D
D=M
@SP
AM=M+1
A=A-1
M=D
@2
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
M=M-D
@0
D=A
@ARG
D=M+D
@R13
M=D
@SP
AM=M-1
D=M
@R13
A=M
M=D
(LOOP)
@0
D=A
@ARG
A=M+D
D=M
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
@COMPUTE_ELEMENT
D;JNE
@END
0;JMP
(COMPUTE_ELEMENT)
@0
D=A
@THAT
A=M+D
D=M
@SP
AM=M+1
A=A-1
M=D
@1
D=A
@THAT
A=M+D
D=M
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
M=M+D
@2
D=A
@THAT
D=M+D
@R13
M=D
@SP
AM=M-1
D=M
@R13
A=M
M=D
@THAT
D=M
@SP
AM=M+1
A=A-1
M=D
@1
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
M=M+D
@SP
AM=M-1
D=M
@THAT
M=D
@0
D=A
@ARG
A=M+D
D=M
@SP
AM=M+1
A=A-1
M=D
@1
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
M=M-D
@0
D=A
@ARG
D=M+D
@R13
M=D
@SP
AM=M-1
D=M
@R13
A=M
M=D
@LOOP
0;JMP
(END)
