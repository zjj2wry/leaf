package main

import (
	"log"

	"github.com/spf13/cobra"
)

func newCommandServer() *cobra.Command {
	leaf := NewLeafCommand()
	var rootCmd = &cobra.Command{Use: "leaf"}
	rootCmd.AddCommand(leaf)
	return rootCmd
}

func main() {
	if err := newCommandServer().Execute(); err != nil {
		log.Fatal(err)
	}
}
