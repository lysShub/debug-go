//go:build debug
// +build debug

// golang zero cost debug-mode simple assertions
package debug

const debug = true

func Debug() bool { return debug }
