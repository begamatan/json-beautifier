package formatter_test

import (
	"errors"
	"testing"

	"github.com/begamatan/json-beautifier/backend/internal/formatter"
)

func TestBeautify(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		indent  int
		want    string
		wantErr error
	}{
		{
			name:   "simple object 2 spaces",
			input:  `{"a":1,"b":2}`,
			indent: 2,
			want:   "{\n  \"a\": 1,\n  \"b\": 2\n}",
		},
		{
			name:   "simple object 4 spaces",
			input:  `{"a":1,"b":2}`,
			indent: 4,
			want:   "{\n    \"a\": 1,\n    \"b\": 2\n}",
		},
		{
			name:   "array",
			input:  `[1,2,3]`,
			indent: 2,
			want:   "[\n  1,\n  2,\n  3\n]",
		},
		{
			name:   "nested object",
			input:  `{"a":{"b":1}}`,
			indent: 2,
			want:   "{\n  \"a\": {\n    \"b\": 1\n  }\n}",
		},
		{
			name:   "already pretty input",
			input:  "{\n  \"a\": 1\n}",
			indent: 2,
			want:   "{\n  \"a\": 1\n}",
		},
		{
			name:    "invalid JSON",
			input:   `{bad}`,
			indent:  2,
			wantErr: formatter.ErrInvalidJSON,
		},
		{
			name:    "empty string",
			input:   ``,
			indent:  2,
			wantErr: formatter.ErrInvalidJSON,
		},
		{
			name:    "invalid indent",
			input:   `{"a":1}`,
			indent:  3,
			wantErr: errors.New("indent must be 2 or 4, got 3"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatter.Beautify([]byte(tc.input), tc.indent)
			if tc.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %q, want %q", string(got), tc.want)
			}
		})
	}
}

func TestMinify(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:  "pretty object",
			input: "{\n  \"a\": 1,\n  \"b\": 2\n}",
			want:  `{"a":1,"b":2}`,
		},
		{
			name:  "already minified",
			input: `{"a":1}`,
			want:  `{"a":1}`,
		},
		{
			name:  "array",
			input: "[\n  1,\n  2\n]",
			want:  `[1,2]`,
		},
		{
			name:    "invalid JSON",
			input:   `{bad}`,
			wantErr: formatter.ErrInvalidJSON,
		},
		{
			name:    "empty input",
			input:   ``,
			wantErr: formatter.ErrInvalidJSON,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatter.Minify([]byte(tc.input))
			if tc.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %q, want %q", string(got), tc.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid object", input: `{"a":1}`, wantErr: false},
		{name: "valid array", input: `[1,2,3]`, wantErr: false},
		{name: "valid string", input: `"hello"`, wantErr: false},
		{name: "valid number", input: `42`, wantErr: false},
		{name: "valid null", input: `null`, wantErr: false},
		{name: "valid bool", input: `true`, wantErr: false},
		{name: "invalid object", input: `{bad}`, wantErr: true},
		{name: "empty", input: ``, wantErr: true},
		{name: "trailing comma", input: `{"a":1,}`, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := formatter.Validate([]byte(tc.input))
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
