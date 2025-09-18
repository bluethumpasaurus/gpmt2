/*
Greenplum Magic Tool

Authored by Tyler Ramer, Ignacio Elizaga, Brian Honohan
Copyright 2018 & 2025

Licensed under the Apache License, Version 2.0 (the "License")
*/
package main

import (
	"fmt"
	"time"

	"github.com/bluethumpasaurus/gpmt2/pkg/db"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os" // Added this import
)

// Tool information constants

const (
	gpmtVersion = "Version (pre)ALPHA"
	githubRepo  = "https://github.com/bluethumpasaurus/gpmt2"
)

type logOptions struct {
	Verbose bool
	LogDir  string
	LogFile string
}

// Local Package Variables
var (
	// gp_log_collector flags
	lcOpts LogCollectorOptions

	// DB connection details
	connString db.ConnString //FIXME/TODO: Do we need a separate wrapper for DB?

	// logging flags
	logOpts = logOptions{LogFile: fmt.Sprintf("/gpmt_log_%s", time.Now().Format("2006-01-02"))}
)

// Sub Command: Version
// When this command is used the version of the gpmt is displayed on the screen
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "GPDB Version number",
	Long:  `Greenplum Magic Tool version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(gpmtVersion)
	},
}

var rootCmd = &cobra.Command{
	Use:   "gpmt",
	Short: "Greenplum Magic Tool 2",
	Long:  `An open-source rewrite of the Greenplum Magic Tool (GPMT) for database diagnostics.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if logOpts.Verbose {
			log.SetLevel(log.DebugLevel)
		}

		// Setup logging to file
		logName := logOpts.LogDir + logOpts.LogFile
		log.SetOutput(os.Stdout)
		formatter := &log.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		}
		log.SetFormatter(formatter)
		log.AddHook(lfshook.NewHook(logName, &log.TextFormatter{}))

	},
	Run: func(cmd *cobra.Command, args []string) {
		// if no argument specified throw the help menu on the screen
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// This is the NEW entry point for the application.
func main() {
	Execute()
}

// Initialize the cobra command CLI.
func init() {
	// All global flags
	rootCmd.PersistentFlags().BoolVarP(&logOpts.Verbose, "verbose", "v", false, "Enable verbose or debug logging")
	rootCmd.PersistentFlags().StringVar(&logOpts.LogDir, "log-directory", "/tmp", "Directory where the logfile should be created") // TODO - logfile default may change

	// Database connection parameters.
	rootCmd.PersistentFlags().StringVar(&connString.Hostname, "hostname", "localhost", "Hostname where the database is hosted")
	rootCmd.PersistentFlags().IntVar(&connString.Port, "port", 5432, "Port number of the master database")
	rootCmd.PersistentFlags().StringVar(&connString.Database, "database", "template1", "Database name to connect")
	rootCmd.PersistentFlags().StringVar(&connString.Username, "username", "gpadmin", "Username that is used to connect to database")
	rootCmd.PersistentFlags().StringVar(&connString.Password, "password", "", "Password for the user")

	// Attach the sub command to the root command.
	rootCmd.AddCommand(versionCmd)

	// NOTE: The other commands (logCollectorCmd, analyzeSessionCmd, etc.)
	// will be added to rootCmd automatically by their own init() functions.
	// Each command should handle its own flag registration in its own init() function.
}
