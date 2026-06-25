package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tsdb",
	Short: "Time series database CLI",
	Long:  `A simple time series database CLI for collecting and querying time series data.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(queryCmd)
}
