package skit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Require(t testing.TB, expected error, actual error, msgAndArgs ...any) {
	if expected == nil {
		require.NoError(t, actual, msgAndArgs...)
	} else {
		require.Error(t, actual, msgAndArgs...)
		require.ErrorIs(t, actual, expected, msgAndArgs...)
	}
}

func RequireNil(t testing.TB, err error, msgAndArgs ...any) {
	require.NoError(t, err, msgAndArgs...)
}

func RequireNotNil(t testing.TB, err error, msgAndArgs ...any) {
	require.Error(t, err, msgAndArgs...)
}
