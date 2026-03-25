package formatter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

// ErrInvalidJSON is returned when the provided input is not valid JSON.
var ErrInvalidJSON = errors.New("invalid JSON")

// Beautify pretty-prints raw JSON with the specified indent size (2 or 4 spaces).
// Returns ErrInvalidJSON if the input cannot be parsed.
func Beautify(input []byte, indent int) ([]byte, error) {
	if indent != 2 && indent != 4 {
		return nil, fmt.Errorf("indent must be 2 or 4, got %d", indent)
	}
	if !json.Valid(input) {
		return nil, ErrInvalidJSON
	}

	var buf bytes.Buffer
	indentStr := spaces(indent)
	if err := json.Indent(&buf, input, "", indentStr); err != nil {
		return nil, ErrInvalidJSON
	}
	return buf.Bytes(), nil
}

// Minify compacts raw JSON by removing unnecessary whitespace.
// Returns ErrInvalidJSON if the input cannot be parsed.
func Minify(input []byte) ([]byte, error) {
	if !json.Valid(input) {
		return nil, ErrInvalidJSON
	}

	var buf bytes.Buffer
	if err := json.Compact(&buf, input); err != nil {
		return nil, ErrInvalidJSON
	}
	return buf.Bytes(), nil
}

// Validate checks whether the input is valid JSON.
// Returns nil if valid, ErrInvalidJSON otherwise.
func Validate(input []byte) error {
	if !json.Valid(input) {
		return ErrInvalidJSON
	}
	return nil
}

func spaces(n int) string {
	s := make([]byte, n)
	for i := range s {
		s[i] = ' '
	}
	return string(s)
}
