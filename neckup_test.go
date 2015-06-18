package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_randomString(test *testing.T) {

	resRandString := randomString(128)

	if len(resRandString) != 128 {
		test.Error("randomString returned invalid length of random string")
	}

}

func Test_uploadHandler(test *testing.T) {

	uploadHandle := http.HandlerFunc(uploadHandler)

	getRequest, _ := http.NewRequest("GET", "", nil)
	getWriter := httptest.NewRecorder()

	uploadHandle.ServeHTTP(getWriter, getRequest)

	if getWriter.Code != http.StatusOK {
		test.Errorf("uploadHandler did not succeed and returned: %v", getWriter.Code)
	}
}

func Test_viewHandler(test *testing.T) {

	writer := httptest.NewRecorder()

	viewHandler(writer, "index", nil)

	if writer.Code != http.StatusOK {
		test.Errorf("viewHandler did not succeed and returned: %v", writer.Code)
	}
}
