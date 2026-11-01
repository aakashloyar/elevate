package process_document

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/ledongthuc/pdf"
)

func extractText(data []byte, fileName, contentType string) (string, error) {
	if isPDF(fileName, contentType) {
		reader, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return "", fmt.Errorf("parse PDF: %w", err)
		}
		textReader, err := reader.GetPlainText()
		if err != nil {
			return "", fmt.Errorf("extract PDF text: %w", err)
		}
		text, err := io.ReadAll(textReader)
		if err != nil {
			return "", fmt.Errorf("read extracted PDF text: %w", err)
		}
		return string(text), nil
	}
	if !utf8.Valid(data) {
		return "", fmt.Errorf("document is not UTF-8 text; unsupported file format")
	}
	return string(data), nil
}

func isPDF(fileName, contentType string) bool {
	return strings.EqualFold(strings.TrimSpace(contentType), "application/pdf") || strings.EqualFold(filepath.Ext(fileName), ".pdf")
}
