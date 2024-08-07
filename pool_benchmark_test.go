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
	"sync"
	"testing"

	"go.mway.dev/pool"
)

func BenchmarkPool(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		pool := sync.Pool{
			New: func() any {
				return newPooledObject()
			},
		}

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				pool.Put(pool.Get())
			}
		})
	})

	b.Run("typed", func(b *testing.B) {
		pool := pool.New(newPooledObject)

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				pool.Put(pool.Get())
			}
		})
	})

	b.Run("typed releaser", func(b *testing.B) {
		pool := pool.NewWithReleaser(newPooledObject, releasePooledObject)

		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				pool.Put(pool.Get())
			}
		})
	})
}

type pooledObject struct {
	a string
	b []byte
	c bool
}

func newPooledObject() *pooledObject {
	return &pooledObject{
		a: "world",
		b: []byte("hello"),
		c: true,
	}
}

func releasePooledObject(o *pooledObject) {
	o.a = ""
	o.b = o.b[:0]
	o.c = false
}
