package main

import (
	"encoding/xml"
	"io"
)

type (
	XMLWriter struct {
		w io.Writer
	}

	XMLToken struct {
		Token
	}
)

func (t XMLToken) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = string(t.Type)
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if t.Value != "" {
		if err := e.EncodeToken(xml.CharData([]byte(t.Value))); err != nil {
			return err
		}
	}
	for _, child := range t.Children {
		encodeChild(e, child)
	}
	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}

	return nil
}

func encodeChild(e *xml.Encoder, child Token) error {
	start := xml.StartElement{Name: xml.Name{Local: string(child.Type)}}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if child.Value != "" {
		if err := e.EncodeToken(xml.CharData([]byte(child.Value))); err != nil {
			return err
		}
	}
	for _, child := range child.Children {
		if err := encodeChild(e, child); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}
	return nil
}

func NewXMLWriter(w io.Writer) XMLWriter {
	return XMLWriter{
		w,
	}
}

func (x XMLWriter) WriteTokens(tokens []Token) error {
	root := XMLToken{
		Token: Token{
			Type:     "tokens",
			Children: tokens,
		},
	}
	output, err := xml.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}
	if _, err := x.w.Write(output); err != nil {
		return err
	}
	return nil
}

func (x XMLWriter) WriteParsedTokens(tokens []Token) error {
	return nil
}
