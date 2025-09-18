# GPMT2 - Greenplum Magic Tool 2

GPMT2 is a Go CLI application for Greenplum database diagnostics and log collection. It provides commands for collecting logs, analyzing sessions, and checking database statistics.

**Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.**

## Working Effectively

### Bootstrap, Build, and Test the Repository
Run these commands in sequence for a complete setup:

1. **Initialize Go modules** (first time only):
   ```bash
   go mod init github.com/bluethumpasaurus/gpmt2
   go mod tidy
   ```
   - `go mod tidy` takes ~13 seconds initially. NEVER CANCEL. Wait for completion.

2. **Build the application**:
   ```bash
   bash scripts/build.sh
   ```
   - Build takes ~0.3 seconds (incremental) or ~13 seconds (clean build). NEVER CANCEL.
   - Creates executable at `build/gpmt`

3. **Format and validate code**:
   ```bash
   go fmt ./...      # ~0.03 seconds
   go vet ./...      # ~0.1 seconds  
   go test ./...     # ~0.1 seconds (no tests currently exist)
   ```

### Run the Application
- **Show help**: `./build/gpmt --help` or `./build/gpmt`
- **Show version**: `./build/gpmt version`
- **Log collector**: `./build/gpmt gp_log_collector --help`

### Application Structure
- Main command: `./build/gpmt` 
- Available subcommands:
  - `version` - Shows application version (currently "Version (pre)ALPHA")
  - `gp_log_collector` - Log collection utility (placeholder implementation)
  - `completion` - Shell completion generation (bash, zsh, fish, powershell)
  - `help` - Help about any command

## Validation

### Always Run These Steps After Making Changes
1. **Format code**: `go fmt ./...` - Must run without errors
2. **Lint code**: `go vet ./...` - Must run without warnings  
3. **Build successfully**: `bash scripts/build.sh` - Must complete without errors
4. **Test basic functionality**:
   ```bash
   ./build/gpmt --help        # Must show help menu
   ./build/gpmt version       # Must show "Version (pre)ALPHA" 
   ./build/gpmt gp_log_collector  # Must show placeholder message
   ```

### Manual Validation Scenarios
After making code changes, always test these scenarios:
- **Clean build**: Remove `build/gpmt go.mod go.sum` and rebuild from scratch
- **Flag handling**: Test `--verbose`, `--hostname`, `--port`, `--database` flags
- **Complex flag combinations**: Test `./build/gpmt --verbose --hostname test --port 5433 gp_log_collector --no-prompts --failed-segs`
- **Error conditions**: Test invalid commands and invalid flags
- **Help system**: Verify `--help` works for all commands
- **Completion**: Test `./build/gpmt completion bash` generates shell completion script

### Complete End-to-End Validation
Run this complete validation sequence after making changes:
```bash
# Format, lint, and build
go fmt ./... && go vet ./... && bash scripts/build.sh

# Test basic functionality
./build/gpmt --help
./build/gpmt version
./build/gpmt gp_log_collector --help

# Test flag combinations
./build/gpmt --verbose --hostname test --port 5433 --database test --username user gp_log_collector --no-prompts

# Test error handling
./build/gpmt invalidcommand  # Should show error and help
./build/gpmt --invalidflag   # Should show error and usage

# Test completion
./build/gpmt completion bash | head -10  # Should generate completion script
```

## Critical Information

### Build Requirements
- **Go version**: 1.24.7+ (confirmed working)
- **Dependencies**: Automatically resolved via `go mod tidy`
- **Build artifacts**: Stored in `build/` directory (gitignored)

### Timing Expectations
- **NEVER CANCEL builds or long-running commands**
- `go mod tidy`: ~13 seconds (first time), ~1 second (subsequent)
- Build: ~0.07 seconds (incremental), ~12 seconds (clean build)  
- Linting: <1 second for all checks (`go fmt` ~0.7s, `go vet` ~3s)
- No timeout issues expected for normal operations

### Known Limitations and Status
- **Current status**: Early development - basic CLI framework implemented
- **gp_log_collector**: Command exists but shows placeholder message only
- **Database features**: Connection parameters accepted but not actively used
- **Testing**: No test files exist yet - `go test` reports "[no test files]"
- **analyze_session and gpstatscheck**: Functions exist but no CLI commands registered

## Common Tasks

### Repository Structure
```
/home/runner/work/gpmt2/gpmt2/
├── build/              # Build artifacts (gitignored)
├── cmd/gpmt/          # Main application source
│   ├── root.go        # CLI root and main()
│   ├── logCollectorCmd.go  # Log collector command definition
│   ├── logCollector.go     # Log collector implementation
│   ├── analyzeSession.go  # Session analysis (placeholder)
│   └── gpstatscheck.go     # Stats checking (placeholder)
├── pkg/db/            # Database connection utilities
│   └── db.go          # PostgreSQL/Greenplum connectivity
├── scripts/           # Build scripts
│   └── build.sh       # Main build script
├── .gitignore         # Git ignore rules
├── go.mod             # Go module definition
├── go.sum             # Go module checksums
└── README.md          # Project documentation
```

### Key Global Flags
All commands support these global database connection flags:
- `--hostname` (default: "localhost")
- `--port` (default: 5432) 
- `--database` (default: "template1")
- `--username` (default: "gpadmin")
- `--password` (default: "")
- `--verbose` / `-v` (enables debug logging)
- `--log-directory` (default: "/tmp")

### Adding New Commands
To add a new command:
1. Create command variable using `cobra.Command`
2. Add `init()` function that calls `rootCmd.AddCommand(yourCmd)`
3. Rebuild with `bash scripts/build.sh`
4. Test with `./build/gpmt --help` to verify registration

### Dependencies and Modules
- Uses Go modules (not dep/vendor)
- Main dependencies: cobra (CLI), logrus (logging), lib/pq (PostgreSQL)
- Run `go mod tidy` after adding new imports
- Module path: `github.com/bluethumpasaurus/gpmt2`
- Legacy files: `Gopkg.toml` and `Gopkg.lock` exist but are ignored (use Go modules)

## Troubleshooting

### Common Issues
1. **"no required module provides package"**: Run `go mod init github.com/bluethumpasaurus/gpmt2 && go mod tidy`
2. **Command not found**: Ensure command has `init()` function calling `rootCmd.AddCommand()`
3. **Build failures**: Check that module path matches in imports and go.mod
4. **Git tracking build artifacts**: Run `git rm --cached build/gpmt` to untrack
5. **Permission denied on executable**: The built binary should be executable: `chmod +x build/gpmt`
6. **Module cache issues**: Clear with `go clean -modcache` then run `go mod tidy`

### Debugging Tips
- Use `--verbose` flag to enable debug logging for any command
- Check build directory exists: `mkdir -p build` if needed
- Legacy files present: `Gopkg.toml` and `Gopkg.lock` exist but are not used (project uses Go modules)

### Environment Variables
- `MASTER_DATA_DIRECTORY`: Used by log collector functionality (currently not active)
- Standard Go environment variables (GOPATH, etc.) work as expected

Remember: This application is in early development. The CLI framework is complete but most business logic shows placeholder messages. Focus changes on the CLI structure, build process, and command registration rather than database-specific functionality.