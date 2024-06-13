package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var (
	url string
	numRequests int
)

func executeRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Println("Response status code:", resp.Status)
}

func run(cmd *cobra.Command, args []string) {
	for i:=0; i<numRequests; i++ {
		executeRequest(url)
	}
}

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "loadtester -u [url] -n [numRequests]",
	Short: "A simple CLI that takes a URL as input",
	Long: "Loadtester is a CLI that can be used to simulate a load on a website or HTTP(S) based API.",
	Run: run,
}

// Execute root command
// This is called by main.main(). 
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// initialize flags
func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL to be requested")
    rootCmd.Flags().IntVarP(&numRequests, "numRequests", "n", 1, "Number of requests to be made")
}


