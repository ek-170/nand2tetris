package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	src := flag.String("src", "", "source file/dir path")
	flag.Parse()

	if src == nil || *src == "" {
		log.Fatal("not set source path")
	}

	if strings.HasSuffix(filepath.Base(*src), jackExt) {
		// process single .jack file
		fmt.Printf("Processing file: %s\n", *src)
		analyzer, err := NewJackAnalyzer(*src)
		defer func() {
			if err := analyzer.Close(); err != nil {
				fmt.Printf("Error while closeing files: %v\n", err)
			}
		}()
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		if err := analyzer.Analyze(); err != nil {
			log.Fatalf("%v\n", err)
		}

	} else {
		err := filepath.WalkDir(*src, func(path string, d os.DirEntry, err error) error {
			// process directory
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if strings.HasSuffix(d.Name(), jackExt) {
				fmt.Printf("Processing file: %s\n", path)
				analyzer, err := NewJackAnalyzer(path)
				if err != nil {
					log.Fatalf("%v\n", err)
				}
				if err := analyzer.Analyze(); err != nil {
					log.Fatalf("%v\n", err)
				}
				if err := analyzer.Close(); err != nil {
					fmt.Printf("Error while closeing files: %v\n", err)
				}
			}
			return nil
		})

		if err != nil {
			fmt.Printf("Error while walking directory: %v\n", err)
			return
		}
	}
}
