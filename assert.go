package debug

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	stdebug "runtime/debug"
	"syscall"
)

func NoError(err error, msg ...any) {
	if err != nil {
		fail(fmt.Sprintf("%+v", err), msg...)
	}
}
func True(v bool, msg ...any) {
	if !v {
		fail("require true", msg...)
	}
}
func False(v bool, msg ...any) {
	if v {
		fail("require false", msg...)
	}
}
func Equal[T comparable](v1, v2 T, msg ...any) {
	if v1 != v2 {
		fail(fmt.Sprintf("require %v == %v", v1, v2), msg...)
	}
}
func NotEqual[T comparable](v1, v2 T, msg ...any) {
	if v1 == v2 {
		fail(fmt.Sprintf("require %v != %v", v1, v2), msg...)
	}
}
func Zero[N number](v N, msg ...any) {
	if v != *new(N) {
		fail(fmt.Sprintf("require zero get %+v", v), msg...)
	}
}
func NotZero[N number](v N, msg ...any) {
	if v == *new(N) {
		fail("require not zero", msg...)
	}
}
func NotEmpty[T any](v T, msg ...any) {
	if empty(v) {
		fail("require not empty", msg...)
	}
}
func Less[N number](v1, v2 N, msg ...any) {
	if v1 >= v2 {
		fail(fmt.Sprintf("require %v < %v", v1, v2), msg...)
	}
}
func Greater[N number](v1, v2 N, msg ...any) {
	if v1 <= v2 {
		fail(fmt.Sprintf("require %v > %v", v1, v2), msg...)
	}
}
func LessOrEqual[N number](v1, v2 N, msg ...any) {
	if v1 > v2 {
		fail(fmt.Sprintf("require %v <= %v", v1, v2), msg...)
	}
}
func GreaterOrEqual[N number](v1, v2 N, msg ...any) {
	if v1 < v2 {
		fail(fmt.Sprintf("require %v >= %v", v1, v2), msg...)
	}
}
func Assert[T any](v any, msg ...any) {
	if _, ok := v.(T); !ok {
		fail(fmt.Sprintf("require %T can assert to %T", v, *new(T)), msg...)
	}
}
func NotAssert[T any](v any, msg ...any) {
	if _, ok := v.(T); ok {
		fail(fmt.Sprintf("require %T cannot assert to %T", v, *new(T)), msg...)
	}
}
func MapHas[K comparable, V any](m map[K]V, key K, msg ...any) {
	_, has := m[key]
	if !has {
		fail(fmt.Sprintf("require map contain key %v", key), msg...)
	}
}
func MapNotHas[K comparable, V any](m map[K]V, key K, msg ...any) {
	val, has := m[key]
	if has {
		fail(fmt.Sprintf("map contain key %v, val %v", key, val), msg...)
	}
}

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func empty(v any) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Chan, reflect.Map, reflect.Slice:
		return val.Len() == 0
	case reflect.Ptr:
		if val.IsNil() {
			return true
		}
		deref := val.Elem().Interface()
		return empty(deref)
	default:
		zero := reflect.Zero(val.Type())
		return reflect.DeepEqual(v, zero.Interface())
	}
}
func fail(s string, msg ...any) {
	var b = &bytes.Buffer{}
	fmt.Fprintln(b, s)
	if len(msg) > 0 {
		fmt.Fprintln(b, "msg:")
		fmt.Fprintln(b, msg...)
	}
	fmt.Fprintln(b, "stack:")
	fmt.Fprintln(b, string(stdebug.Stack()))

	Fail(b.String())
}

var Fail = func(s string) {
	fmt.Fprintln(os.Stderr, s)
	os.Exit(int(syscall.SIGABRT))
}
