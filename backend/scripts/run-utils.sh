#!/bin/bash

# Space Notifications Backend Utility Runner
# This script provides an easy way to run utility commands

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Function to check if required environment variables are set
check_env() {
    local missing=()
    
    if [ -z "$PB_URL" ]; then
        missing+=("PB_URL")
    fi
    
    if [ -z "$PB_ADMIN_EMAIL" ]; then
        missing+=("PB_ADMIN_EMAIL")
    fi
    
    if [ -z "$PB_ADMIN_PASSWORD" ]; then
        missing+=("PB_ADMIN_PASSWORD")
    fi
    
    if [ ${#missing[@]} -ne 0 ]; then
        print_error "Missing required environment variables: ${missing[*]}"
        echo ""
        echo "Please set the following environment variables:"
        echo "  export PB_URL=\"http://localhost:8080\""
        echo "  export PB_ADMIN_EMAIL=\"your-admin@email.com\""
        echo "  export PB_ADMIN_PASSWORD=\"your-admin-password\""
        echo ""
        exit 1
    fi
}

# Function to show help
show_help() {
    echo "Space Notifications Backend Utility Runner"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Available commands:"
    echo "  cleanup-events    Remove duplicate events from the database"
    echo "  help             Show this help message"
    echo ""
    echo "Environment variables required:"
    echo "  PB_URL           PocketBase URL (e.g., http://localhost:8080)"
    echo "  PB_ADMIN_EMAIL   PocketBase admin email"
    echo "  PB_ADMIN_PASSWORD PocketBase admin password"
    echo ""
    echo "Examples:"
    echo "  $0 cleanup-events"
    echo "  $0 help"
}

# Change to the backend directory (assuming script is in backend/scripts/)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"
cd "$BACKEND_DIR"

# Main command handling
case "${1:-help}" in
    "cleanup-events")
        print_info "Running duplicate events cleanup..."
        check_env
        go run cmd/utils-main.go -cleanup-events
        print_success "Cleanup utility completed!"
        ;;
    "help"|"--help"|"-h")
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac
