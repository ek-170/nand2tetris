package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	src := flag.String("src", "", "source file/dir path")
	dest := flag.String("dest", "", "output file path")
	flag.Parse()

	if src == nil || *src == "" {
		fmt.Println("not set source path")
		return
	}

	if dest == nil || *dest == "" {
		fmt.Println("not set output path")
		return
	}

	p := NewParserWithFile(*dest)
	defer func() {
		if err := p.Close(); err != nil {
			fmt.Printf("Error while closeing files: %v\n", err)
		}
	}()

	if strings.HasSuffix(filepath.Base(*src), ".vm") {
		fmt.Printf("Processing file: %s\n", *src)
		file, err := os.Open(*src)
		if err != nil {
			fmt.Printf("Failed to open file %s: %v\n", *src, err)
		}
		p.Do(file)
	} else {
		err := filepath.WalkDir(*src, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if strings.HasSuffix(d.Name(), ".vm") {
				fmt.Printf("Processing file: %s\n", path)
				file, err := os.Open(path)
				if err != nil {
					fmt.Printf("Failed to open file %s: %v\n", path, err)
					return nil
				}
				p.Do(file)
			}
			return nil
		})

		if err != nil {
			fmt.Printf("Error while walking directory: %v\n", err)
			return
		}
	}
}
