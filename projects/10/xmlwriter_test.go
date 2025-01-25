package main

import (
	"os"
	"testing"
)

func TestXMLWriter(t *testing.T) {
	tokens := []Token{
		{
			Type: "keyword",
			Children: []Token{
				{
					Type:  "keyword",
					Value: "while",
					Children: []Token{
						{
							Type:  "simbol",
							Value: "+",
						},
					},
				},
				{
					Type:  "simbol",
					Value: ";",
				},
			},
		},
		{
			Type:  "simbol",
			Value: "{",
		},
		{
			Type:  "simbol",
			Value: "}",
		},
	}
	writer := NewXMLWriter(os.Stdout)
	if err := writer.WriteTokens(tokens); err != nil {
		t.Fatal(err)
	}
}
