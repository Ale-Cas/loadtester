package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCLIApplication(t *testing.T) {
	// Start a localhost server to simulate the API
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	// Redirect stdout to a buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Set the arguments
	nReq := 2
	RootCmd.SetArgs([]string{"-u", server.URL, "-n", fmt.Sprint(nReq)})

	// Run the CLI
	Execute()

	// Stop capturing stdout
	w.Close()
	os.Stdout = old
	// Read the buffer
	var buf bytes.Buffer
	buf.ReadFrom(r)

	// Check the output
	outPerReq := "Response status code: 200 OK\n"
	expected := ""
	for i := 0; i < nReq; i++ {
		expected += outPerReq
	}
	if buf.String() != expected {
		t.Errorf("Expected %s, got %s", expected, buf.String())
	}
}
