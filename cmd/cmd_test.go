package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCli(t *testing.T) {
	// Simulate 2 requests to a server
	// 1 success and 1 failure (not found)
	var reqs int
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		reqs++
		if reqs > 1 {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusOK)
		}
	}))

	defer server.Close()

	// Set the arguments
	nReq := 2
	RootCmd.SetArgs([]string{"-u", server.URL, "-n", fmt.Sprint(nReq)})

	// Run the CLI
	output := captureStdout(func() {
		Execute()
	})

	// Check the output
	expected := "Successes: 1\nFailures: 1\nErrors: 0\n"
	if output != expected {
		t.Errorf("Expected %s, got %s", expected, output)
	}
}

func TestCliDebugMode(t *testing.T) {
	// Start a localhost server to simulate the API
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	// Set the arguments
	nReq := 2
	RootCmd.SetArgs([]string{"-u", server.URL, "-n", fmt.Sprint(nReq), "-d"})

	// Run the CLI
	output := captureStdout(func() {
		Execute()
	})

	// Check the output
	outPerReq := "Response Status: 200 OK\n"
	expected := ""
	for i := 0; i < nReq; i++ {
		expected += outPerReq
	}
	expected += fmt.Sprintf("Successes: %d\nFailures: 0\nErrors: 0\n", nReq)
	if output != expected {
		t.Errorf("Expected %s, got %s", expected, output)
	}
}
