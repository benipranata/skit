package skit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Require asserts that actual matches expected error.
func Require(t testing.TB, expected error, actual error, msgAndArgs ...any) {
	if expected == nil {
		require.NoError(t, actual, msgAndArgs...)
	} else {
		require.Error(t, actual, msgAndArgs...)
		require.ErrorIs(t, actual, expected, msgAndArgs...)
	}
}

// RequireNil asserts that err is nil.
func RequireNil(t testing.TB, err error, msgAndArgs ...any) {
	require.NoError(t, err, msgAndArgs...)
}

// RequireNotNil asserts that err is not nil.
func RequireNotNil(t testing.TB, err error, msgAndArgs ...any) {
	require.Error(t, err, msgAndArgs...)
}
