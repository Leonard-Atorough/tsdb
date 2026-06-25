package main

import (
	"time"

	"github.com/spf13/cobra"
)

var (
	dataDir         string
	tenantID        string
	pollingInterval time.Duration
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start collecting time series data",
	Long:  `Start the collector to continuously gather time series data and write it to files.`,
	Run: func(cmd *cobra.Command, args []string) {
		start(dataDir, tenantID, pollingInterval)
	},
}

func init() {
	startCmd.Flags().StringVar(&dataDir, "datadir", "data", "Data directory path")
	startCmd.Flags().StringVar(&tenantID, "tenant", "tenant1", "Tenant ID")
	startCmd.Flags().DurationVar(&pollingInterval, "interval", 500*time.Millisecond, "Polling interval")
}
