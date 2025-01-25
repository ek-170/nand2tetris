package main

import (
	"fmt"
	"os"
	"testing"
)

func TestTokenize(t *testing.T) {
	srcFile, err := os.Open("./ArrayTest/Main.jack")
	if err != nil {
		t.Fatal(err)
	}
	tokens, err := NewJackTokenizer(srcFile).Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tokens)
}
