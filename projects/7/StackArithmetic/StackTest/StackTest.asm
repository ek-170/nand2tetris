// ---./StackArithmetic/StackTest/StackTest.vm---
@17
D=A
@SP
AM=M+1
A=A-1
M=D
@17
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@END_EQ1
D;JEQ
@SP
A=M-1
M=0
(END_EQ1)
@17
D=A
@SP
AM=M+1
A=A-1
M=D
@16
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@END_EQ2
D;JEQ
@SP
A=M-1
M=0
(END_EQ2)
@16
D=A
@SP
AM=M+1
A=A-1
M=D
@17
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@END_EQ3
D;JEQ
@SP
A=M-1
M=0
(END_EQ3)
@892
D=A
@SP
AM=M+1
A=A-1
M=D
@891
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@END_LT1
D;JLT
@SP
A=M-1
M=0
(END_LT1)
@891
D=A
@SP
AM=M+1
A=A-1
M=D
@892
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@END_LT2
D;JLT
@SP
A=M-1
M=0
(END_LT2)
@891
D=A
@SP
AM=M+1
A=A-1
M=D
@891
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@END_LT3
D;JLT
@SP
A=M-1
M=0
(END_LT3)
@32767
D=A
@SP
AM=M+1
A=A-1
M=D
@32766
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@END_GT1
D;JGT
@SP
A=M-1
M=0
(END_GT1)
@32766
D=A
@SP
AM=M+1
A=A-1
M=D
@32767
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@END_GT2
D;JGT
@SP
A=M-1
M=0
(END_GT2)
@32766
D=A
@SP
AM=M+1
A=A-1
M=D
@32766
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@END_GT3
D;JGT
@SP
A=M-1
M=0
(END_GT3)
@57
D=A
@SP
AM=M+1
A=A-1
M=D
@31
D=A
@SP
AM=M+1
A=A-1
M=D
@53
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
@112
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
@SP
A=M-1
M=-M
@SP
AM=M-1
D=M
A=A-1
M=M&D
@82
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
M=M|D
@SP
A=M-1
M=!M
