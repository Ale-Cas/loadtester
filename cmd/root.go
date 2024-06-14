package cmd

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var (
	url           string
	numRequests   int  // total num of requests to be made
	numConcurrent int  // num of concurrent requests to be made
	debug         bool // debug mode, if true print all the responses
)


// executeRequest makes a GET request to the given URL
// and returns the response status code or an error
func executeRequest(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		if debug {
			fmt.Println("Error:", err)
		}
		return 0, err
	}
	defer resp.Body.Close()
	if debug {
		fmt.Println("Response Status:", resp.Status)
	}
	return resp.StatusCode, nil
}

func run(cmd *cobra.Command, args []string) {
	// counters for the responses
	var successes, failures, errorCount int
	var wg sync.WaitGroup
	// semaphore to limit concurrent goroutines
	semaphore := make(chan bool, numConcurrent)
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			semaphore <- true
			status, err := executeRequest(url)
			// >= 200 and < 300 is considered a success
			if err != nil {
				errorCount++
			} else if status >= http.StatusOK && status < http.StatusMultipleChoices {
				successes++
			} else {
				failures++
			}
			<-semaphore
			wg.Done()
		}()
	}
	wg.Wait() // wait for all goroutines to finish
	
	// print the results
	fmt.Println("Successes:", successes)
	fmt.Println("Failures:", failures)
	fmt.Println("Errors:", errorCount)
}

// RootCmd represents the base command
var RootCmd = &cobra.Command{
	Use:   "loadtester -u [url] -n [numRequests]",
	Short: "A simple CLI that takes a URL as input",
	Long:  "Loadtester is a CLI that can be used to simulate a load on a website or HTTP(S) based API.",
	Run:   run,
}

// Execute root command
// This is called by main.main().
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// initialize flags
func init() {
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.Flags().StringVarP(&url, "url", "u", "", "URL to be requested")
	RootCmd.Flags().IntVarP(&numRequests, "numRequests", "n", 1, "Total number of requests to be made")
	RootCmd.Flags().IntVarP(&numConcurrent, "numConcurrent", "c", 1, "Number of concurrent requests to be made")
	RootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug mode if true print all the responses")
}
