// Copyright © 2018 Xavier Basty <xbasty@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"go/scanner"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
)

// Assert fails the test if the condition is false.
// Taken from https://github.com/benbjohnson/testing.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	tb.Helper()
	if !condition {
		tb.Logf(msg, v...)
		tb.FailNow()
	}
}

// Ok fails the test if an err is not nil.
// Taken from https://github.com/benbjohnson/testing.
func Ok(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Logf("unexpected error: %s", err.Error())
		tb.FailNow()
	}
}

// Equals fails the test if exp is not equal to act.
// Taken from https://github.com/benbjohnson/testing.
func Equals(tb testing.TB, exp, act interface{}) {
	tb.Helper()
	if diff := deep.Equal(exp, act); diff != nil {
		tb.Logf("%s\nexp: %sgot: %s", diff, spew.Sdump(exp), spew.Sdump(act))
		tb.FailNow()
	}
}

// ErrEquals fails the test if act is nil or act.Error() != exp
func ErrEquals(tb testing.TB, exp string, act error) {
	tb.Helper()
	if act == nil {
		tb.Logf("exp err %q but err was nil\n", exp)
		tb.FailNow()
	}
	if act.Error() != exp {
		tb.Logf("exp err: %q but got: %q\n", exp, act.Error())
		tb.FailNow()
	}
}

func ScannerErrEquals(tb testing.TB, exp string, act error) {
	tb.Helper()
	if act == nil {
		tb.Logf("exp err %q but err was nil\n", "exp")
		tb.FailNow()
	}
	switch e := errors.Cause(act).(type) {
	case scanner.ErrorList:
		if e[0].Error() != exp {
			tb.Logf("exp err: %q but got: %q\n", exp, e[0].Error())
			tb.FailNow()
		}
	case scanner.Error:
		if e.Error() != "foo" {
			tb.Logf("exp err: %q but got: %q\n", "exp", e.Error())
			tb.FailNow()
		}
	default:
		tb.Logf("exp scanner err but got: %T\n", e)
	}
}
