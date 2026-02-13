package handlers

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockNetError struct {
	message string
	timeout bool
}

func (e mockNetError) Error() string {
	return e.message
}

func (e mockNetError) Timeout() bool {
	return e.timeout
}

func (e mockNetError) Temporary() bool {
	return false
}

func TestCategorizeError(t *testing.T) {
	plainErr := errors.New("plain error")

	tests := []struct {
		name     string
		input    error
		expected string
		nilErr   bool
	}{
		{
			name:   "nil error",
			input:  nil,
			nilErr: true,
		},
		{
			name:     "timeout network error",
			input:    mockNetError{message: "dial timeout", timeout: true},
			expected: "TIMEOUT: dial timeout",
		},
		{
			name:     "non-timeout network error",
			input:    mockNetError{message: "connection refused", timeout: false},
			expected: "NETWORK_ERROR: connection refused",
		},
		{
			name:     "eof error",
			input:    io.EOF,
			expected: "CONNECTION_CLOSED: Remote host closed connection",
		},
		{
			name:     "plain error passthrough",
			input:    plainErr,
			expected: "plain error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := categorizeError(tt.input)
			if tt.nilErr {
				assert.Nil(t, got)
				return
			}

			assert.Equal(t, tt.expected, got.Error())
		})
	}
}

func TestResultToString(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectedLen int
		expected    []string
	}{
		{
			name:        "nil input",
			input:       nil,
			expectedLen: 0,
		},
		{
			name:        "bool slice",
			input:       []bool{true, false, true},
			expectedLen: 3,
		},
		{
			name:        "uint16 slice",
			input:       []uint16{1, 2},
			expectedLen: 2,
		},
		{
			name:        "int slice",
			input:       []int{3, 4, 5, 6},
			expectedLen: 4,
		},
		{
			name:        "interface slice",
			input:       []interface{}{1, "x"},
			expectedLen: 2,
		},
		{
			name:        "string slice passthrough",
			input:       []string{"a", "b"},
			expectedLen: 2,
			expected:    []string{"a", "b"},
		},
		{
			name:        "unknown type",
			input:       struct{ Value int }{Value: 1},
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := resultToString(tt.input)
			assert.Len(t, got, tt.expectedLen)
			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
