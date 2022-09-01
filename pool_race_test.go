// Copyright (c) 2022 Matt Way
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE THE SOFTWARE.

package pool_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/pool"
)

type pooledValue struct {
	n int
}

func TestPool_Race(t *testing.T) {
	var (
		n int
		p = pool.New(func() *pooledValue {
			n++
			return &pooledValue{
				n: n,
			}
		})
	)

	for i := 0; i < 1_000; i++ {
		x := p.Get()
		require.Equal(t, i+1, x.n)
	}
}

func TestPool_Releaser_Race(t *testing.T) {
	var (
		p = pool.NewWithReleaser(
			func() *pooledValue {
				return nil
			},
			func(x *pooledValue) {
				x.n = -1
			},
		)
		x = &pooledValue{
			n: 123,
		}
	)

	p.Put(x)
	require.Equal(t, -1, x.n)
}
