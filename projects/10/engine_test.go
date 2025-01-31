package main

import (
	"os"
	"strings"
	"testing"
)

const (
	src1 = `class Main {
		field int x, y; 
		field int size; 
	}
	`
	src2 = `class Main {
    static boolean test;    // Added for testing -- there is no static keyword
                            // in the Square files.

    function void main() {
        var SquareGame game;
    }

    function void more() {  // Added to test Jack syntax that is not used in
        var boolean b;      // the Square files.
    }
}
		`
)

func TestCompile(t *testing.T) {
	tokens, err := NewJackTokenizer(strings.NewReader(src2)).Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	engine := NewCompilationEngine(tokens)
	parsed, err := engine.Parse()
	if err != nil {
		t.Fatal(err)
	}
	writer := NewXMLWriter(os.Stdout)
	if err := writer.WriteParsedTokens(parsed); err != nil {
		t.Fatal(err)
	}
}
