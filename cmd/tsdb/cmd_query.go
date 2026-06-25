package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/leonard-atorough/tsdb/internal"
	"github.com/spf13/cobra"
)

var (
	startTime string
	endTime   string
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query time series data",
	Long:  `Query time series data within a time range.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := &internal.Config{
			DataDir:  dataDir,
			TenantID: tenantID,
		}
		filePath := config.GetFilePath("data")
		reader := internal.NewReader(filePath)
		
		results, err := reader.Query(startTime, endTime)
		if err != nil {
			log.Fatalf("Error querying data: %v", err)
		}
		
		// Output results as JSON
		fmt.Fprintf(log.Writer(), "Found %d results\n", len(results))
		for _, result := range results {
			jsonBytes, err := json.Marshal(result)
			if err != nil {
				log.Printf("Error marshaling result: %v", err)
				continue
			}
			fmt.Println(string(jsonBytes))
		}
	},
}

func init() {
	queryCmd.Flags().StringVar(&startTime, "start", "", "Start time (RFC3339 format)")
	queryCmd.Flags().StringVar(&endTime, "end", "", "End time (RFC3339 format)")
	queryCmd.Flags().StringVar(&dataDir, "datadir", "data", "Data directory path")
	queryCmd.Flags().StringVar(&tenantID, "tenant", "tenant1", "Tenant ID")
	queryCmd.MarkFlagRequired("start")
	queryCmd.MarkFlagRequired("end")
	queryCmd.MarkFlagRequired("datadir")
	queryCmd.MarkFlagRequired("tenant")
}
