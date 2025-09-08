#!/bin/bash
set -eu

# ==================================================================================== #
# VARIABLES
# ==================================================================================== #

# Load environment variables from .env file if it exists
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Check if DSN is set
if [ -z "${DSN:-}" ]; then
    echo "Error: DSN environment variable is required"
    echo "Either set it in your .env file or export it:"
    echo "export DSN='postgres://user:password@localhost/database?sslmode=disable'"
    exit 1
fi

# ==================================================================================== #
# SCRIPT LOGIC
# ==================================================================================== #

# Parse command line arguments
NAME=""
EXPIRES_DAYS=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--name)
            NAME="$2"
            shift 2
            ;;
        -e|--expires)
            EXPIRES_DAYS="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 -n|--name <key_name> [-e|--expires <days>]"
            echo ""
            echo "Options:"
            echo "  -n, --name      API key name/description (required)"
            echo "  -e, --expires   Expiry in days (optional, 0 = no expiry)"
            echo "  -h, --help      Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0 --name \"Production API Key\""
            echo "  $0 --name \"Temp Key\" --expires 30"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done

# Validate required parameters
if [ -z "$NAME" ]; then
    echo "Error: API key name is required"
    echo "Use: $0 --name \"Your Key Name\""
    echo "Use -h or --help for more options"
    exit 1
fi

echo "Creating API key: $NAME"

# Generate random API key (64 hex characters)
API_KEY=$(openssl rand -hex 32)

# Hash the key for storage
KEY_HASH=$(echo -n "$API_KEY" | shasum -a 256 | cut -d' ' -f1)

# Prepare expiry date if specified
EXPIRES_CLAUSE=""
if [ -n "$EXPIRES_DAYS" ] && [ "$EXPIRES_DAYS" -gt 0 ]; then
    EXPIRES_DATE=$(date -d "+${EXPIRES_DAYS} days" '+%Y-%m-%d %H:%M:%S')
    EXPIRES_CLAUSE="'$EXPIRES_DATE'"
    echo "Key will expire on: $EXPIRES_DATE"
else
    EXPIRES_CLAUSE="NULL"
    echo "Key will never expire"
fi

# Insert into database
QUERY="
INSERT INTO api_keys (name, key_hash, is_active, created_at, updated_at, expires_at)
VALUES ('$NAME', '$KEY_HASH', true, NOW(), NOW(), $EXPIRES_CLAUSE)
RETURNING id, created_at;
"

# Execute query and capture results
RESULT=$(psql "$DSN" -t -c "$QUERY" 2>/dev/null) || {
    echo "Error: Failed to insert API key into database"
    echo "Make sure your database is running and DSN is correct"
    exit 1
}

# Parse the result
ID=$(echo "$RESULT" | awk '{print $1}' | tr -d ' ')
CREATED_AT=$(echo "$RESULT" | awk '{print $3, $4}' | tr -d ' ')

echo ""
echo "✅ API Key created successfully!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "ID:      $ID"
echo "Name:    $NAME"
echo "Created: $CREATED_AT"
echo ""
echo "🔑 API Key (save this securely - it won't be shown again):"
echo "$API_KEY"
echo ""
echo "📝 Usage example:"
echo "Authorization: Bearer $API_KEY"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"