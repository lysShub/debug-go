package debug_test

import (
	"errors"
	"math"
	"os"
	"testing"

	"github.com/lysShub/debug-go"
)

func Test_Assert(t *testing.T) {
	back := debug.Fail
	var called int
	debug.Fail = func(string) { called++ }
	defer func() { debug.Fail = back }()

	for _, fn := range []func(){
		func() { debug.NoError(errors.New("111")) },
		func() { debug.True(false) },
		func() { debug.False(true) },
		func() { debug.Zero(1) },
		func() { debug.Zero(math.SmallestNonzeroFloat64) },
		func() { debug.NotZero(0) },
		func() { debug.NotZero(0.0) },
		func() { debug.NotEmpty("") },
		func() { debug.NotEmpty(make([]byte, 0)) },
		func() { debug.NotEmpty(make(chan int)) },
		func() { debug.Less(1, 1) },
		func() { debug.Less(2, 1) },
		func() { debug.Greater(1, 1) },
		func() { debug.Greater(1, 2) },
		func() { debug.LessOrEqual(2, 1) },
		func() { debug.GreaterOrEqual(1, 2) },
		func() { debug.Assert[error]("ssss") },
		func() { debug.NotAssert[error](os.ErrClosed) },
		func() { debug.MapHas(map[int]int{0: 0}, 1) },
		func() { debug.MapNotHas(map[int]int{0: 0}, 0) },
	} {
		old := called
		fn()

		if old+1 != called {
			t.Fatalf("%v  %v", old+1, called)
		}
	}
}
