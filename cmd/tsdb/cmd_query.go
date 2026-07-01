package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/leonard-atorough/tsdb/internal"
	"github.com/leonard-atorough/tsdb/internal/reader"
	"github.com/spf13/cobra"
)

var (
	startTime   string
	endTime     string
	agg         []string
	field       string
	measurement string
	tags        []string
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
		r, err := reader.NewReader(filePath)
		if err != nil {
			log.Fatalf("Error creating reader: %v", err)
		}

		if agg != nil {
			if field == "" {
				log.Fatalf("Field name is required for aggregation")
			}
			if measurement == "" {
				log.Fatalf("Measurement name is required for aggregation")
			}

			tagsMap := makeTagsMap(tags)

			opts := reader.AggregateOpts{
				Field:       field,
				Measurement: measurement,
				Funcs:       agg,
				From:        startTime,
				To:          endTime,
				Tags:        tagsMap,
			}
			result, err := r.Aggregates(opts)
			if err != nil {
				log.Fatalf("Error performing aggregation: %v", err)
			}

			fmt.Fprintf(log.Writer(), "Aggregation results for measurement '%s' field '%s':\n", measurement, field)
			for aggFunc, value := range result {
				fmt.Fprintf(log.Writer(), "%s: %v\n", aggFunc, value)
			}
		} else {
			tagsMap := makeTagsMap(tags)

			opts := reader.QueryOpts{
				From:        startTime,
				To:          endTime,
				Tags:        tagsMap,
				Measurement: measurement,
			}
			result, err := r.Query(opts)
			if err != nil {
				log.Fatalf("Error querying data: %v", err)
			}
			// Output results as JSON
			fmt.Fprintf(log.Writer(), "Found %d results\n", len(result))
			for _, r := range result {
				jsonBytes, err := json.Marshal(r)
				if err != nil {
					log.Printf("Error marshaling result: %v", err)
					continue
				}
				fmt.Println(string(jsonBytes))
			}
		}

	},
}

func init() {
	queryCmd.Flags().StringVarP(&dataDir, "datadir", "d", "data", "Data directory path")
	queryCmd.Flags().StringVarP(&tenantID, "tenant", "t", "tenant1", "Tenant ID")
	queryCmd.Flags().StringVarP(&startTime, "start", "s", "", "Start time (RFC3339 format)")
	queryCmd.Flags().StringVarP(&endTime, "end", "e", "", "End time (RFC3339 format)")
	queryCmd.Flags().StringVarP(&measurement, "measurement", "m", "", "Measurement name for aggregation")
	queryCmd.Flags().StringVarP(&field, "field", "f", "", "Field name for aggregation")
	queryCmd.Flags().StringArrayVarP(&tags, "tags", "t", []string{}, "Tags for filtering")
	queryCmd.Flags().StringArrayVarP(&agg, "agg", "a", []string{}, "Aggregation function (avg, sum, count, min, max)")
	queryCmd.MarkFlagRequired("datadir")
	queryCmd.MarkFlagRequired("tenant")
}

func makeTagsMap(tags []string) map[string]string {
	tagsMap := make(map[string]string)
	for _, tag := range tags {
		parts := strings.SplitN(tag, "=", 2)
		if len(parts) != 2 {
			log.Fatalf("Invalid tag format: %s. Expected format: key=value", tag)
		}
		tagsMap[parts[0]] = parts[1]
	}
	return tagsMap
}
