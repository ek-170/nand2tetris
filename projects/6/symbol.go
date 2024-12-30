package main

import "fmt"

type SymbolTable map[string]int

var ramAddress = 16
var symbolTable = newSymbolTable()

func newSymbolTable() SymbolTable {
	s := map[string]int{
		"SP":     0x0000,
		"LCL":    0x0001,
		"ARG":    0x0002,
		"THIS":   0x0003,
		"THAT":   0x0004,
		"R0":     0x0000,
		"R1":     0x0001,
		"R2":     0x0002,
		"R3":     0x0003,
		"R4":     0x0004,
		"R5":     0x0005,
		"R6":     0x0006,
		"R7":     0x0007,
		"R8":     0x0008,
		"R9":     0x0009,
		"R10":    0x000A,
		"R11":    0x000B,
		"R12":    0x000C,
		"R13":    0x000D,
		"R14":    0x000E,
		"R15":    0x000F,
		"SCREEN": 0x4000,
		"KBD":    0x6000,
	}
	return s
}

func (s SymbolTable) AddRAMEntry(symbol string) {
	if s.Contains(symbol) {
		// fmt.Printf("symbol %v has already set, update symbol's address: %v\n", symbol, ramAddress)
		fmt.Printf("symbol %v has already set with RAM address: %v\n", symbol, symbolTable[symbol])
		return
	}
	s[symbol] = ramAddress
	ramAddress++
}

func (s SymbolTable) AddROMEntry(symbol string, address int) {
	if s.Contains(symbol) {
		// fmt.Printf("symbol %v has already set, update symbol's address: %v\n", symbol, address)
		fmt.Printf("symbol %v has already set with ROM address: %v\n", symbol, address)
		return
	}
	s[symbol] = address
}

func (s SymbolTable) Contains(symbol string) bool {
	_, ok := s[symbol]
	return ok
}

func (s SymbolTable) GetAddress(symbol string) (address int, ok bool) {
	address, ok = s[symbol]
	return address, ok
}
