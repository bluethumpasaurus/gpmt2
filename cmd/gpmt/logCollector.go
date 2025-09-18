package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// getLogDirectoryFromDB queries the database to get the actual log directory path
func getLogDirectoryFromDB() (string, error) {
	const query = "select distinct datadir || '/log' from gp_segment_configuration where content='-1';"

	log.Debug("Querying database for log directory path")

	// Try to execute the query, but handle any panics from connection failures
	var result []map[string]interface{}
	var err error

	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("database connection failed: %v", r)
			}
		}()
		result, err = connString.ExecuteQuery(query)
	}()

	if err != nil {
		return "", fmt.Errorf("failed to query database for log directory: %w", err)
	}

	if len(result) == 0 {
		return "", fmt.Errorf("no log directory found in gp_segment_configuration")
	}

	// Extract the directory path from the result
	for _, row := range result {
		for _, value := range row {
			if str, ok := value.(string); ok && str != "" {
				logDir := strings.TrimSpace(str)
				log.Debugf("Found log directory from database: %s", logDir)
				return logDir, nil
			} else if bytes, ok := value.([]byte); ok && len(bytes) > 0 {
				logDir := strings.TrimSpace(string(bytes))
				log.Debugf("Found log directory from database: %s", logDir)
				return logDir, nil
			}
		}
	}

	return "", fmt.Errorf("invalid log directory result from database")
}

// logCollector archives Greenplum Database log files from the master and segment directories.
func logCollector(archiveName string) error {
	// Default to a timestamped archive name if none is provided.
	if archiveName == "" {
		timestamp := time.Now().Format("20060102_150405")
		archiveName = fmt.Sprintf("gpmt_logs_%s.tar.gz", timestamp)
	}

	fmt.Printf("Starting log collection...\n")
	fmt.Printf("Logs will be archived to: %s\n", archiveName)

	// Create the output file.
	outFile, err := os.Create(archiveName)
	if err != nil {
		return fmt.Errorf("failed to create archive file: %w", err)
	}
	defer outFile.Close()

	// Create a gzip writer.
	gw := gzip.NewWriter(outFile)
	defer gw.Close()

	// Create a tar writer.
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Get the log directory from the database first
	logDir, err := getLogDirectoryFromDB()
	if err != nil {
		log.Debugf("Failed to get log directory from database: %v", err)

		// Fallback to environment variable or hardcoded path
		gpMasterDir := os.Getenv("MASTER_DATA_DIRECTORY")
		if gpMasterDir == "" {
			// Fallback for when the environment variable is not set.
			homeDir, err := os.UserHomeDir()
			if err == nil {
				gpMasterDir = filepath.Join(homeDir, "gpdb", "gp-master", "gpseg-1")
			}
		}

		if gpMasterDir == "" {
			return fmt.Errorf("unable to determine log directory: database query failed and MASTER_DATA_DIRECTORY environment variable not set")
		}

		logDir = filepath.Join(gpMasterDir, "pg_log")
		log.Debugf("Using fallback log directory: %s", logDir)
	}

	// Walk the log directory and add files to the archive.
	err = filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Add the file to the tar archive.
		return addFileToTar(tw, path, logDir)
	})

	if err != nil {
		return fmt.Errorf("failed to walk log directory: %w", err)
	}

	fmt.Println("Log collection complete.")
	return nil
}

// addFileToTar is a helper function to add a file to a tar archive.
func addFileToTar(tw *tar.Writer, path string, basePath string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(stat, stat.Name())
	if err != nil {
		return err
	}

	// Use a relative path in the archive.
	// First try MASTER_DATA_DIRECTORY, then fall back to basePath
	masterDataDir := os.Getenv("MASTER_DATA_DIRECTORY")
	if masterDataDir != "" {
		header.Name = strings.TrimPrefix(path, masterDataDir)
	} else if basePath != "" {
		header.Name = strings.TrimPrefix(path, basePath)
	} else {
		header.Name = filepath.Base(path)
	}

	// Ensure header name doesn't start with /
	header.Name = strings.TrimPrefix(header.Name, "/")

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tw, file); err != nil {
		return err
	}

	fmt.Printf("  - Archived %s\n", path)
	return nil
}
