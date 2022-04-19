package pool

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPoolUnconvertiblePanic(t *testing.T) {
	cases := []struct {
		input       any
		expectPanic bool
	}{
		{
			input:       int64(123),
			expectPanic: false,
		},
		{
			input:       int32(123),
			expectPanic: true,
		},
		{
			input:       uint64(123),
			expectPanic: true,
		},
		{
			input:       int(123),
			expectPanic: true,
		},
		{
			input:       uint(123),
			expectPanic: true,
		},
		{
			input:       nil,
			expectPanic: true,
		},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("%T", tt.input), func(t *testing.T) {
			pool := New(func() int64 { return 0 })
			pool.pool.New = func() any {
				return tt.input
			}

			fn := require.Panics
			if !tt.expectPanic {
				fn = require.NotPanics
			}

			fn(t, func() {
				pool.pool.Put(tt.input)
				pool.Get()
			})
		})
	}
}
