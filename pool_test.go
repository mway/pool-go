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

//go:build !race

package pool_test

import (
	"bytes"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/pool"
)

func TestPool_Constructor(t *testing.T) {
	defer debug.SetGCPercent(debug.SetGCPercent(-1))

	p := pool.New(func() *bytes.Buffer {
		return bytes.NewBuffer([]byte(t.Name()))
	})

	for i := 0; i < 1_000; i++ {
		func() {
			buf := p.Get()
			defer p.Put(buf)
			require.Equal(t, t.Name(), buf.String())
		}()
	}
}

func TestPool_Releaser(t *testing.T) {
	defer debug.SetGCPercent(debug.SetGCPercent(-1))

	p := pool.NewWithReleaser(
		func() *bytes.Buffer {
			return bytes.NewBuffer([]byte(t.Name()))
		},
		func(x *bytes.Buffer) {
			x.Reset()
		},
	)

	tmp := make([]*bytes.Buffer, 1_000)
	for i := 0; i < len(tmp); i++ {
		tmp[i] = p.Get()
	}
	for i := 0; i < len(tmp); i++ {
		p.Put(tmp[i])
	}

	for i := 0; i < 1_000; i++ {
		buf := p.Get()
		require.Equal(t, 0, buf.Len(), buf.String())
	}
}
