package internal

import (
	"path/filepath"
	"time"
)

type Config struct {
	DataDir         string        `json:"data_dir"`
	TenantID        string        `json:"tenant_id"`
	PollingInterval time.Duration `json:"polling_interval"`
}

func (c *Config) GetFilePath(metric string) string {
	return filepath.Join(c.DataDir, c.TenantID, metric+".jsonl") // Construct the file path based on the data directory, tenant ID, and metric name
}
