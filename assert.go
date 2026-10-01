package debug

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	stdebug "runtime/debug"
	"syscall"
)

func NoError(err error, msgAndArgs ...any) {
	if err != nil {
		fail(fmt.Sprintf("%+v", err), msgAndArgs...)
	}
}
func True(v bool, msgAndArgs ...any) {
	if !v {
		fail("require true", msgAndArgs...)
	}
}
func False(v bool, msgAndArgs ...any) {
	if v {
		fail("require false", msgAndArgs...)
	}
}
func Equal[T comparable](v1, v2 T, msgAndArgs ...any) {
	if v1 != v2 {
		fail(fmt.Sprintf("require %v == %v", v1, v2), msgAndArgs...)
	}
}
func NotEqual[T comparable](v1, v2 T, msgAndArgs ...any) {
	if v1 == v2 {
		fail(fmt.Sprintf("require %v != %v", v1, v2), msgAndArgs...)
	}
}
func Zero[N number](v N, msgAndArgs ...any) {
	if v != *new(N) {
		fail(fmt.Sprintf("require zero get %+v", v), msgAndArgs...)
	}
}
func NotZero[N number](v N, msgAndArgs ...any) {
	if v == *new(N) {
		fail("require not zero", msgAndArgs...)
	}
}
func NotEmpty[T any](v T, msgAndArgs ...any) {
	if empty(v) {
		fail("require not empty", msgAndArgs...)
	}
}
func Less[N number](v1, v2 N, msgAndArgs ...any) {
	if v1 >= v2 {
		fail(fmt.Sprintf("require %v < %v", v1, v2), msgAndArgs...)
	}
}
func Greater[N number](v1, v2 N, msgAndArgs ...any) {
	if v1 <= v2 {
		fail(fmt.Sprintf("require %v > %v", v1, v2), msgAndArgs...)
	}
}
func LessOrEqual[N number](v1, v2 N, msgAndArgs ...any) {
	if v1 > v2 {
		fail(fmt.Sprintf("require %v <= %v", v1, v2), msgAndArgs...)
	}
}
func GreaterOrEqual[N number](v1, v2 N, msgAndArgs ...any) {
	if v1 < v2 {
		fail(fmt.Sprintf("require %v >= %v", v1, v2), msgAndArgs...)
	}
}
func Assert[T any](v any, msgAndArgs ...any) {
	if _, ok := v.(T); !ok {
		fail(fmt.Sprintf("require %T can assert to %T", v, *new(T)), msgAndArgs...)
	}
}
func NotAssert[T any](v any, msgAndArgs ...any) {
	if _, ok := v.(T); ok {
		fail(fmt.Sprintf("require %T cannot assert to %T", v, *new(T)), msgAndArgs...)
	}
}
func MapHas[K comparable, V any](m map[K]V, key K, msgAndArgs ...any) {
	_, has := m[key]
	if !has {
		fail(fmt.Sprintf("require map contain key %v", key), msgAndArgs...)
	}
}
func MapNotHas[K comparable, V any](m map[K]V, key K, msgAndArgs ...any) {
	val, has := m[key]
	if has {
		fail(fmt.Sprintf("map contain key %v, val %v", key, val), msgAndArgs...)
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
func fail(s string, msgAndArgs ...any) {
	var b = &bytes.Buffer{}
	fmt.Fprintln(b, s)
	if len(msgAndArgs) > 0 {
		fmt.Fprintln(b, "msg:")
		b.WriteString(messageFromMsgAndArgs(msgAndArgs...))
	}
	fmt.Fprintln(b, "stack:")
	fmt.Fprintln(b, string(stdebug.Stack()))

	Fail(b.String())
}
func messageFromMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msg)
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}

var Fail = func(s string) {
	fmt.Fprintln(os.Stderr, s)
	os.Exit(int(syscall.SIGABRT))
}
