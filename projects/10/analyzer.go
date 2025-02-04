package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const jackExt = ".jack"

type JackWriter interface {
	WriteTokens(tokens []*Token) error
	WriteParsedTokens(tokens *Token) error
}

type JackAnalyzer struct {
	src, dst  *os.File
	tokenizer JackTokenizer
	w         JackWriter
}

func NewJackAnalyzer(srcFilePath string) (JackAnalyzer, error) {
	if exists := ExistsFilePath(srcFilePath); !exists {
		return JackAnalyzer{}, fmt.Errorf("expect file path, but %q is directory: ", srcFilePath)
	}
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return JackAnalyzer{}, err
	}
	distFileName := strings.TrimSuffix(filepath.Base(srcFilePath), jackExt)
	dstDir := filepath.Dir(srcFilePath)
	dstPath := filepath.Join(dstDir, distFileName)
	// 既存のxmlの削除防止のために_をつける
	dstFile, err := OpenFileWithReset(dstPath + "_p.xml")
	if err != nil {
		return JackAnalyzer{}, fmt.Errorf("failed to open or create file: %w", err)
	}
	return JackAnalyzer{
		src:       srcFile,
		dst:       dstFile,
		tokenizer: NewJackTokenizer(srcFile),
		w:         NewXMLWriter(dstFile),
	}, nil
}

func (ja JackAnalyzer) Analyze() error {
	tokens, err := ja.tokenizer.Tokenize()
	if err != nil {
		return nil
	}
	// if err := ja.w.WriteTokens(tokens); err != nil {
	// 	return err
	// }
	engine := NewCompilationEngine(tokens)
	parsed, err := engine.Parse()
	if err != nil {
		return err
	}
	if err := ja.w.WriteParsedTokens(parsed); err != nil {
		return err
	}
	return nil
}

func (ja JackAnalyzer) Close() error {
	var err error

	if ja.src != nil {
		if closeErr := ja.src.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close reader: %w", closeErr)
		}
	}

	if ja.dst != nil {
		if closeErr := ja.dst.Close(); closeErr != nil {
			if err != nil {
				err = fmt.Errorf("%v; failed to close writer: %w", err, closeErr)
			} else {
				err = fmt.Errorf("failed to close writer: %w", closeErr)
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
}
