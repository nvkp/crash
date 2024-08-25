package crash_test

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func equal(t *testing.T, expected, actual interface{}, msg string, args ...interface{}) {
	if eq(expected, actual) {
		return
	}

	t.Errorf("%s: expected: %s\n actual: %s", fmt.Sprintf(msg, args...), expected, actual)
}

func errorIs(t *testing.T, err, target error, msg string, args ...interface{}) {
	if errors.Is(err, target) {
		return
	}

	t.Errorf("%s: error: %v is not error: %v", fmt.Sprintf(msg, args...), err, target)
}

func eq(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}
