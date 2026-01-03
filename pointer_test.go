package skit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPtr(t *testing.T) {
	sampleString := "test"

	tests := []struct {
		name string
		v    any
	}{
		{name: "string", v: "test"},
		{name: "error", v: ErrNew("test")},
		{name: "struct", v: struct{ Name string }{Name: "test"}},
		{name: "nil", v: nil},
		{name: "pointer", v: &sampleString},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Ptr(tt.v)
			require.NotNil(t, got)
			require.Equal(t, tt.v, *got)
		})
	}
}

func TestDeref(t *testing.T) {
	sampleString := "test"
	sampleErr := ErrNew("test")
	type sample struct{ Name string }
	sampleStruct := sample{Name: "test"}

	tests := []struct {
		name string
		fn   func() any
		want any
	}{
		{name: "string", fn: func() any { return Deref(&sampleString) }, want: "test"},
		{name: "error", fn: func() any { return Deref(&sampleErr) }, want: sampleErr},
		{name: "struct", fn: func() any { return Deref(&sampleStruct) }, want: sampleStruct},
		{name: "nil", fn: func() any { return Deref((*string)(nil)) }, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.fn())
		})
	}
}

func TestDerefOr(t *testing.T) {
	sampleString := "test"

	tests := []struct {
		name string
		fn   func() any
		want any
	}{
		{name: "string", fn: func() any { return DerefOr(&sampleString, "fallback") }, want: "test"},
		{name: "nil string", fn: func() any { return DerefOr((*string)(nil), "fallback") }, want: "fallback"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.fn())
		})
	}
}
