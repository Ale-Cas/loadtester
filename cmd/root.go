package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var (
	url         string
	numRequests int
)

func executeRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Response status code:", resp.Status)
}

func run(cmd *cobra.Command, args []string) {
	for i := 0; i < numRequests; i++ {
		executeRequest(url)
	}
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
	RootCmd.Flags().IntVarP(&numRequests, "numRequests", "n", 1, "Number of requests to be made")
}
