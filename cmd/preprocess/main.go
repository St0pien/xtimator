package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scrapper [id]",
	Short: "fetch X replies to tweet with specified [id]",
	Args:  cobra.MinimumNArgs(1),
	Run:   PreProcess,
}

func main() {
	var apiKey string

	rootCmd.Flags().StringVar(&apiKey, "api-key", os.Getenv("RAPID_API_KEY"), "Rapid API key, if empty scrapper will pull api key from environmental variable RAPID_API_KEY")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
