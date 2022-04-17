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
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mway.dev/pool"
)

type pooledObject struct {
	a []byte
	b string
	c bool
}

func newPooledObject() *pooledObject {
	return &pooledObject{
		a: make([]byte, 0, 128),
		b: strconv.Itoa(int(time.Now().UnixNano())),
		c: true,
	}
}

func releasePooledObject(o *pooledObject) {
	o.a = o.a[:0]
	o.b = ""
	o.c = false
}

func TestPool(t *testing.T) {
	var (
		pool    = pool.New(newPooledObject)
		expectB = "hello, world!"
		expectA = []byte(expectB)
		expectC = true
	)

	ob := pool.Get()
	ob.a = expectA
	ob.b = expectB
	ob.c = expectC
	pool.Put(ob)

	ob = pool.Get()
	require.Equal(t, expectA, ob.a, "expected object not de-pooled")
	require.Equal(t, expectB, ob.b, "expected object not de-pooled")
	require.Equal(t, expectC, ob.c, "expected object not de-pooled")
}

func TestPoolReleaser(t *testing.T) {
	const iterations = 1000

	var (
		pool = pool.NewWithReleaser(newPooledObject, releasePooledObject)
		wg   sync.WaitGroup
	)

	for i := 0; i < iterations; i++ {
		ob := newPooledObject()
		ob.b = "foo"
		ob.a = []byte(ob.b)
		ob.c = true
		pool.Put(ob)
	}

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ob := pool.Get()
			require.Len(t, ob.a, 0)
			require.Len(t, ob.b, 0)
			require.False(t, ob.c)

			if t.Failed() {
				require.FailNow(t, "releaser did not clear pooled object: %+v", ob)
			}

			ob.b = strconv.Itoa(int(time.Now().UnixNano()))
			ob.a = []byte(ob.b)
			ob.c = true

			pool.Put(ob)
		}()
	}

	wg.Wait()
}
