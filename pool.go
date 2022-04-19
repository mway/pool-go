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

// Package pool provides types and functionality related to object pooling.
package pool

import (
	"errors"
	"fmt"
	"sync"
)

// ErrCorruptPool is returned by Pool.Get when the de-pooled object's type is
// unconvertible to the pool's intended object type.
var ErrCorruptPool = errors.New("corrupt pool")

type (
	// A Constructor creates a new object of type T.
	Constructor[T any] func() T
	// A Releaser resets an object of type T for reuse.
	Releaser[T any] func(T)
)

// Pool is a strongly-typed object pool.
type Pool[T any] struct {
	pool    sync.Pool
	release Releaser[T]
}

// New creates a new Pool compatible with objects of type T with the given
// Constructor for T.
func New[T any](ctor Constructor[T]) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return ctor()
			},
		},
	}
}

// NewWithReleaser creates a new Pool compatible with objects of type T with
// the given Constructor and Releaser for T.
func NewWithReleaser[T any](ctor Constructor[T], release Releaser[T]) *Pool[T] {
	pool := New(ctor)
	pool.release = release
	return pool
}

// Get de-pools or creates an object of type T.
func (p *Pool[T]) Get() T {
	x, ok := p.pool.Get().(T)
	if !ok {
		panic(fmt.Errorf("%v: pool contains non-%T", ErrCorruptPool, x))
	}

	return x
}

// Put places the given object back into the pool. If a Releaser is configured,
// it will be called prior to pooling the object.
func (p *Pool[T]) Put(x T) {
	if p.release != nil {
		p.release(x)
	}

	p.pool.Put(x)
}
