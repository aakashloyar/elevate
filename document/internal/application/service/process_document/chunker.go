package process_document

import (
	"strings"
	"unicode/utf8"
)

const (
	defaultChunkSize    = 1200
	defaultChunkOverlap = 150
)

func splitText(text string) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}

	// Keep paragraph boundaries where possible, while using a bounded rune
	// window so a very long paragraph cannot create an oversized embedding.
	runes := []rune(text)
	chunks := make([]string, 0, (len(runes)/defaultChunkSize)+1)
	for start := 0; start < len(runes); {
		end := start + defaultChunkSize
		if end > len(runes) {
			end = len(runes)
		}
		if end < len(runes) {
			if boundary := strings.LastIndexAny(string(runes[start:end]), "\n.!?"); boundary > defaultChunkSize/2 {
				end = start + boundary + 1
			}
		}
		chunk := strings.TrimSpace(string(runes[start:end]))
		if chunk != "" {
			chunks = append(chunks, chunk)
		}
		if end == len(runes) {
			break
		}
		next := end - defaultChunkOverlap
		if next <= start {
			next = end
		}
		start = next
	}
	return chunks
}

func validText(data []byte) bool {
	return utf8.Valid(data)
}
