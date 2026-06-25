package internal

import (
	"log"
	"os"
	"path/filepath"

	"github.com/leonard-atorough/tsdb/internal/models"
)

type FileWriter struct {
	file *os.File
	path string
}



// NewFileWriter creates a new FileWriter instance for the specified file path.
func NewFileWriter(path string) *FileWriter {
	projectRoot, err := getProjectRoot()
	if err != nil {
		log.Fatalf("Error finding project root: %v", err)
	}
	fullPath := filepath.Join(projectRoot, path)
	err = os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating directory for file: %v", err)
	}

	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	return &FileWriter{
		file: f,
		path: fullPath,
	}
}

// WriteData writes a TimeSeriesData instance to the file in JSON format, appending a newline.
func (fw *FileWriter) WriteData(data models.TimeSeriesData) error {
	jsonData, err := models.MarshalLine(&data)
	if err != nil {
		log.Fatalf("Error marshalling data: %v", err)
		return err
	}
	_, err = fw.file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing data to file: %v", err)
		return err
	}
	return nil
}

// Close closes the file writer, ensuring that all data is flushed to disk.
func (fw *FileWriter) Close() error {
	err := fw.file.Close()
	if err != nil {
		log.Fatalf("Error closing file: %v", err)
		return err
	}
	return nil
}