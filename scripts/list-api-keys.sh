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
SHOW_INACTIVE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -a|--all)
            SHOW_INACTIVE=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [-a|--all]"
            echo ""
            echo "Options:"
            echo "  -a, --all       Show inactive keys as well"
            echo "  -h, --help      Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done

# Build WHERE clause
WHERE_CLAUSE=""
if [ "$SHOW_INACTIVE" = false ]; then
    WHERE_CLAUSE="WHERE is_active = true"
fi

# Query to get API keys
QUERY="
SELECT 
    id,
    name,
    is_active,
    created_at,
    last_used_at,
    expires_at,
    CASE 
        WHEN expires_at IS NOT NULL AND expires_at < NOW() THEN true
        ELSE false
    END as is_expired
FROM api_keys 
$WHERE_CLAUSE
ORDER BY created_at DESC;
"

echo "API Keys:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Execute query and format output
psql "$DSN" -c "$QUERY" --no-align --field-separator='|' --tuples-only 2>/dev/null | while IFS='|' read -r id name is_active created_at last_used_at expires_at is_expired; do
    
    # Format active status
    if [ "$is_active" = "t" ]; then
        if [ "$is_expired" = "t" ]; then
            status="🔴 EXPIRED"
        else
            status="🟢 Active"
        fi
    else
        status="⚫ Inactive"
    fi
    
    # Format last used
    if [ "$last_used_at" = "" ] || [ "$last_used_at" = " " ]; then
        last_used="Never"
    else
        last_used=$(date -d "$last_used_at" '+%Y-%m-%d %H:%M' 2>/dev/null || echo "$last_used_at")
    fi
    
    # Format expires
    if [ "$expires_at" = "" ] || [ "$expires_at" = " " ]; then
        expires="Never"
    else
        expires=$(date -d "$expires_at" '+%Y-%m-%d %H:%M' 2>/dev/null || echo "$expires_at")
    fi
    
    # Format created
    created=$(date -d "$created_at" '+%Y-%m-%d %H:%M' 2>/dev/null || echo "$created_at")
    
    # Truncate name if too long
    if [ ${#name} -gt 25 ]; then
        display_name="${name:0:22}..."
    else
        display_name="$name"
    fi
    
    printf "ID: %-5s | %-28s | %s\n" "$id" "$display_name" "$status"
    printf "         Created: %-16s | Last Used: %-16s | Expires: %s\n" "$created" "$last_used" "$expires"
    echo "         ────────────────────────────────────────────────────────────────────────────────"
    
done || {
    echo "Error: Failed to query API keys"
    echo "Make sure your database is running and DSN is correct"
    exit 1
}

echo ""
if [ "$SHOW_INACTIVE" = false ]; then
    echo "💡 Use --all to show inactive keys as well"
fi