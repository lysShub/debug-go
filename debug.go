//go:build debug
// +build debug

// golang zero cost debug-mode simple assertions
package debug

const debug = true

//go:noinline
func Debug() bool { return debug }
