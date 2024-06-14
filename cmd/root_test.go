package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func captureStdout(f func()) string {
	// Keep a backup of the real stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function that will write to stdout
	f()

	// Stop capturing stdout
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)

	return buf.String()
}

func testExecuteRequest(t *testing.T, status int, expectedOutput string) {
	debug = true // enable debug mode to print the response
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(status)
	}))

	defer server.Close()

	output := captureStdout(func() {
		executeRequest(server.URL)
	})

	// Check the output
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected '%s' in output, got %s", expectedOutput, output)
	}
}

func TestExecuteRequest_OK(t *testing.T) {
	testExecuteRequest(t, http.StatusOK, "200 OK")
}

func TestExecuteRequest_NotFound(t *testing.T) {
	testExecuteRequest(t, http.StatusNotFound, "404 Not Found")
}

// TestExecuteRequestError tests the executeRequest function with an invalid URL
func TestExecuteRequestError(t *testing.T) {
	debug = true
	output := captureStdout(func() {
		executeRequest("invalid")
	})

	// Check the output
	if !strings.Contains(output, "unsupported protocol scheme") {
		t.Errorf("Expected 'unsupported protocol scheme' in output, got %s", output)
	}

	output = captureStdout(func() {
		executeRequest("http://urldoesntexist/")
	})

	// Check the output
	if !strings.Contains(output, "no such host") {
		t.Errorf("Expected 'no such host' in output, got %s", output)
	}
}