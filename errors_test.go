package skit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrNew(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{name: "simple error", text: "something went wrong"},
		{name: "empty error", text: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrNew(tt.text)
			if tt.text == "" {
				RequireNil(t, err)
			} else {
				RequireNotNil(t, err)
				require.EqualError(t, err, tt.text)
			}
		})
	}
}

func TestErrFormat(t *testing.T) {
	tests := []struct {
		name   string
		format string
		args   []any
		exp    string
	}{
		{
			name:   "no args",
			format: "simple error",
			args:   nil,
			exp:    "simple error",
		},
		{
			name:   "with string",
			format: "error: %s",
			args:   []any{"details"},
			exp:    "error: details",
		},
		{
			name:   "with int",
			format: "code: %d",
			args:   []any{42},
			exp:    "code: 42",
		},
		{
			name:   "multiple args",
			format: "%s failed with code %d",
			args:   []any{"operation", 500},
			exp:    "operation failed with code 500",
		},
		{
			name:   "empty format",
			format: "",
			args:   nil,
			exp:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrFormat(tt.format, tt.args...)
			if tt.format == "" {
				RequireNil(t, err)
			} else {
				RequireNotNil(t, err)
				require.EqualError(t, err, tt.exp)
			}
		})
	}
}

func TestErrWrap(t *testing.T) {
	baseErr := errors.New("base error")

	tests := []struct {
		name string
		err  error
		text string
		exp  string
	}{
		{
			name: "nil error",
			err:  nil,
			text: "context",
		},
		{
			name: "simple wrap",
			err:  baseErr,
			text: "context",
			exp:  "context: base error",
		},
		{
			name: "empty text",
			err:  baseErr,
			text: "",
			exp:  "base error",
		},
		{
			name: "whitespace text",
			err:  baseErr,
			text: "   ",
			exp:  "base error",
		},
		{
			name: "text with percent",
			err:  baseErr,
			text: "100% complete",
			exp:  "100%% complete: base error",
		},
		{
			name: "trimmed text",
			err:  baseErr,
			text: "  padded  ",
			exp:  "padded: base error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrWrap(tt.err, tt.text)
			Require(t, tt.err, err)
			if err != nil {
				require.EqualError(t, err, tt.exp)
			}
		})
	}
}
