package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// Helper functions from: https://github.com/benbjohnson/testing
// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func Test_randomString(test *testing.T) {
	const randLength = 128

	equals(test, randLength, len(randomString(randLength)))
}

func Test_uploadHandler(test *testing.T) {

	uploadHandle := http.HandlerFunc(uploadHandler)

	getRequest, err := http.NewRequest("GET", "", nil)
	ok(test, err)

	getWriter := httptest.NewRecorder()

	uploadHandle.ServeHTTP(getWriter, getRequest)
	equals(test, getWriter.Code, http.StatusOK)
}

func Test_viewHandler(test *testing.T) {

	writer := httptest.NewRecorder()

	viewHandler(writer, "index", nil)
	equals(test, writer.Code, http.StatusOK)
}
