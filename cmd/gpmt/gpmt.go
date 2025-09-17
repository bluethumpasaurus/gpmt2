/* Greenplum magic tool

Authored by Tyler Ramer, Ignacio Elizaga
Copyright 2018

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package main

import (
   "archive/tar"
   "compress/gzip"
   "flag"
   "fmt"
   "io"
   "os"
   "path/filepath"
   "strings"
   "time"
   log "github.com/sirupsen/logrus"
)

func init() {
   log.SetOutput(os.Stdout)
}

// DBNAME is the default database name to check for in logs.
const DBNAME = "gpadmin"

// AppName is the name of the application.
const AppName = "gpmt"


func main() {
   // Define subcommands using flagsets.
   logCollectorCmd := flag.NewFlagSet("gp_log_collector", flag.ExitOnError)
   analyzeSessionCmd := flag.NewFlagSet("analyze_session", flag.ExitOnError)
   gpstatscheckCmd := flag.NewFlagSet("gpstatscheck", flag.ExitOnError)

   // Add a flag for the log collector's output file.
   archiveName := logCollectorCmd.String("o", "", "Output archive name (e.g., my_logs.tar.gz)")

   // Check for the correct number of arguments.
   if len(os.Args) < 2 {
      fmt.Printf("Usage: %s <command> [options]\n", AppName)
      fmt.Println("\nAvailable commands:")
      fmt.Println("  gp_log_collector  Collect Greenplum Database log files.")
      fmt.Println("  analyze_session   Analyze active and recent database sessions.")
      fmt.Println("  gpstatscheck      Check for missing or stale table statistics.")
      return
   }

   // Parse the subcommand.
   switch os.Args[1] {
   case "gp_log_collector":
      logCollectorCmd.Parse(os.Args[2:])
      if err := logCollector(*archiveName); err != nil {
         fmt.Fprintf(os.Stderr, "Error: %v\n", err)
         os.Exit(1)
      }
   case "analyze_session":
      analyzeSessionCmd.Parse(os.Args[2:])
      if err := analyzeSession(); err != nil {
         fmt.Fprintf(os.Stderr, "Error: %v\n", err)
         os.Exit(1)
      }
   case "gpstatscheck":
      gpstatscheckCmd.Parse(os.Args[2:])
      if err := gpstatscheck(); err != nil {
         fmt.Fprintf(os.Stderr, "Error: %v\n", err)
         os.Exit(1)
      }
   default:
      fmt.Printf("Unknown command: %s\n", os.Args[1])
      flag.PrintDefaults()
      os.Exit(1)
   }
}