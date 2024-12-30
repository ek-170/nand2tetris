package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	src := flag.String("src", "", "assembly file path")
	dest := flag.String("dest", "", "hack file path")
	flag.Parse()

	if src == nil || *src == "" {
		fmt.Println("not set assembly file path")
		return
	}

	if dest == nil || *dest == "" {
		fmt.Println("not set hack file name")
		return
	}

	f, err := os.Open(*src)
	if err != nil {
		fmt.Println("could not open assembly file")
		return
	}
	p := NewParserWithFile(f, *dest)
	p.Do()
}
