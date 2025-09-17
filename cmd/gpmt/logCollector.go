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
)

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

	// Locate the Greenplum data directories.
	gpMasterDir := os.Getenv("MASTER_DATA_DIRECTORY")
	if gpMasterDir == "" {
		// Fallback for when the environment variable is not set.
		homeDir, err := os.UserHomeDir()
		if err == nil {
			gpMasterDir = filepath.Join(homeDir, "gpdb", "gp-master", "gpseg-1")
		}
	}

	if gpMasterDir == "" {
		return fmt.Errorf("MASTER_DATA_DIRECTORY environment variable not set")
	}

	logDir := filepath.Join(gpMasterDir, "pg_log")

	// Walk the log directory and add files to the archive.
	err = filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Add the file to the tar archive.
		return addFileToTar(tw, path)
	})

	if err != nil {
		return fmt.Errorf("failed to walk log directory: %w", err)
	}

	fmt.Println("Log collection complete.")
	return nil
}

// addFileToTar is a helper function to add a file to a tar archive.
func addFileToTar(tw *tar.Writer, path string) error {
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
	header.Name = strings.TrimPrefix(path, os.Getenv("MASTER_DATA_DIRECTORY"))

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tw, file); err != nil {
		return err
	}

	fmt.Printf("  - Archived %s\n", path)
	return nil
}