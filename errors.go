package skit

import (
	"errors"
	"fmt"
	"strings"
)

// ErrNew creates a new error with the given text.
// This is a simple wrapper around [errors.New].
func ErrNew(text string) error {
	return errors.New(text)
}

// ErrFormat creates a new formatted error with the given format and arguments.
// This is a simple wrapper around [fmt.Errorf].
func ErrFormat(textFormat string, args ...interface{}) error {
	return fmt.Errorf(textFormat, args...)
}

// ErrWrap wraps an existing error with additional context text.
func ErrWrap(err error, text string) error {
	if err == nil {
		return nil
	}
	trimmedText := strings.TrimSpace(text)
	if trimmedText == "" {
		return err
	}
	return fmt.Errorf("%s: %w", strings.ReplaceAll(trimmedText, "%", "%%"), err)
}
