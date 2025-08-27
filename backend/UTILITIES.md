# Space Notifications Backend Utilities

This document describes the utility commands available for the Space Notifications backend system.

## Overview

The backend includes several utility commands that can be run independently of the main backend application. These utilities help with maintenance tasks like cleaning up duplicate data.

## Prerequisites

Before running any utility commands, ensure you have the required environment variables set:

```bash
export PB_URL="http://localhost:8080"              # PocketBase URL
export PB_ADMIN_EMAIL="your-admin@email.com"       # PocketBase admin email
export PB_ADMIN_PASSWORD="your-admin-password"     # PocketBase admin password
```

## Available Utilities

### Cleanup Events

Removes duplicate events from the database based on normalized event titles.

**What it does:**
- Fetches all events from the PocketBase database
- Identifies events with duplicate titles (case-insensitive, trimmed)
- Keeps the first occurrence and removes subsequent duplicates
- Logs all operations for transparency

## Running Utilities

There are three different ways to run the utility commands:

### Method 1: Using the Makefile (Recommended)

```bash
# Show available commands
make help

# Run cleanup events utility
make utils-cleanup-events-dev

# Show utility help
make utils-help-dev

# Build and run utilities (creates binaries)
make build-utils
make utils-cleanup-events
```

### Method 2: Using the Shell Script

```bash
# Show help
./scripts/run-utils.sh help

# Run cleanup events
./scripts/run-utils.sh cleanup-events
```

### Method 3: Direct Go Commands

```bash
# Show utility help
go run cmd/utils-main.go -help

# Run cleanup events
go run cmd/utils-main.go -cleanup-events
```

## Safety Features

- **Environment Validation**: All utilities validate that required environment variables are set before proceeding
- **Connection Retry**: Built-in retry logic for PocketBase connections (similar to the main backend)
- **Detailed Logging**: All operations are logged with clear success/failure messages
- **Non-Disruptive**: Utilities can be run while the main backend is running without interference

## Adding New Utilities

To add a new utility function:

1. Create the function in the appropriate package under `internal/utils/`
2. Add a new command-line flag in `cmd/utils-main.go`
3. Add the flag handling logic in the main function
4. Update the help text and this documentation
5. Add corresponding Makefile targets and shell script commands

Example:
```go
// In cmd/utils-main.go
var (
    cleanupEvents = flag.Bool("cleanup-events", false, "Remove duplicate events from the database")
    newUtility    = flag.Bool("new-utility", false, "Description of new utility")
)

// In the main function
if *newUtility {
    log.Println("Starting new utility...")
    if err := utils.NewUtilityFunction(client); err != nil {
        log.Fatalf("New utility failed: %v", err)
    }
    log.Println("✅ New utility completed successfully!")
}
```

## Troubleshooting

### Common Issues

1. **"Missing environment variables" error**
   - Ensure all required environment variables are set
   - Check that PocketBase is running on the specified URL

2. **"Admin login failed" error**
   - Verify admin credentials are correct
   - Ensure PocketBase is accessible and ready
   - The utility includes retry logic, so transient connection issues should resolve automatically

3. **"No action specified" error**
   - You must specify at least one utility flag (e.g., `-cleanup-events`)
   - Use `-help` to see available options

### Logs

All utilities provide detailed logging including:
- Connection status and retry attempts
- Number of records processed
- Specific actions taken (e.g., which duplicates were removed)
- Success/failure status

## Development Notes

- Utilities use the same configuration system as the main backend
- The retry logic is reduced (10 attempts vs 30) for faster failure feedback
- All utilities exit cleanly with appropriate status codes
- The code is structured to make adding new utilities straightforward
