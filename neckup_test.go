package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
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

const randLength = 24

func Test_randomString(test *testing.T) {
	equals(test, len(randomString(randLength)), randLength)
}

func Test_uploadHandler_GET(test *testing.T) {

	handle := uploadHandler()

	req, err := http.NewRequest("GET", "", nil)
	rec := httptest.NewRecorder()
	ok(test, err)

	handle.ServeHTTP(rec, req)
	equals(test, http.StatusOK, rec.Code) // uploadHandler should call viewHandler which compiles the "index" view
}

func Test_uploadHandler_POST(test *testing.T) {

	// Find "file.mock" (multi-line regexp)
	reg := regexp.MustCompile("(?m)file.mock")

	// Open mock file
	path := "./mocks/file.mock"
	file, err := os.Open(path)
	ok(test, err)
	defer file.Close()

	// Initialize handle, body and multipart writer
	handle := uploadHandler()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create form in body with one file
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	ok(test, err)

	// Copy file into part
	_, err = io.Copy(part, file)
	ok(test, err)

	// Close writer
	err = writer.Close()
	ok(test, err)

	// Create request and start httptest recorder
	req, err := http.NewRequest("POST", "/", body)
	req.Header.Add("Content-Type", writer.FormDataContentType()) // This little fucker got me stuck for hours
	rec := httptest.NewRecorder()
	ok(test, err)

	// Call handle through serveHTTP using recorder and request
	handle.ServeHTTP(rec, req)
	equals(test, http.StatusOK, rec.Code)
	equals(test, true, reg.MatchString(rec.Body.String())) // Test if viewHandler compile view with the filename

}

func Test_viewHandler(test *testing.T) {

	rec := httptest.NewRecorder()

	viewHandler(rec, "minimal", nil)
	equals(test, http.StatusOK, rec.Code) // viewHandler should compile views without file map
}

func Benchmark_randomString(bench *testing.B) {

	for n := 0; n < bench.N; n++ {
		randomString(randLength)
	}
}
