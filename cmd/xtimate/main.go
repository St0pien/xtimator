package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "xtimate [path]",
	Short: "Xtimate conntation of responses in [path]",
	Args:  cobra.MinimumNArgs(1),
	Run:   Xtimate,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
