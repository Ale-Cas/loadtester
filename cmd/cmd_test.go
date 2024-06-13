package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCLIApplication(t *testing.T) {
	// Start a localhost server to simulate the API
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
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
	outPerReq := "Response status code: 200 OK\n"
	expected := ""
	for i := 0; i < nReq; i++ {
		expected += outPerReq
	}
	if output != expected {
		t.Errorf("Expected %s, got %s", expected, output)
	}
}
