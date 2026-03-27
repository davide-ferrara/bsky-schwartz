#!/bin/bash

# ==========================================
# Download Bluesky feeds by Schwartz values
# Creates JSON files ready for /api/analysis/by-url
# ==========================================

set -e

API_URL="http://localhost:8080/api/search"
MODEL="gpt4o"
LIMIT=50  # ← Change this to adjust number of results per query

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUTPUT_DIR="$SCRIPT_DIR/feeds"

# Ensure output directory exists
mkdir -p "$OUTPUT_DIR"

# Function to fetch URLs and accumulate them
fetch_urls() {
    local query="$1"
    
    echo "  Fetching: $query"
    
    # Fetch URLs from API
    urls=$(curl -s "${API_URL}?query=${query}&limit=${LIMIT}" | jq -r '.[]' 2>/dev/null)
    
    if [ -z "$urls" ]; then
        echo "    ⚠️  No results"
        return
    fi
    
    # Count URLs
    url_count=$(echo "$urls" | wc -l)
    echo "    ✓ Found $url_count URLs"
    
    # Return URLs
    echo "$urls"
}

# Function to create feed JSON from multiple queries
create_feed() {
    local group_name="$1"
    shift
    local queries=("$@")
    
    echo ""
    echo "=========================================="
    echo "Downloading: $group_name"
    echo "=========================================="
    
    all_urls=""
    
    for query in "${queries[@]}"; do
        urls=$(fetch_urls "$query")
        if [ -n "$urls" ]; then
            all_urls="$all_urls$urls"$'\n'
        fi
    done
    
    if [ -z "$all_urls" ]; then
        echo "  ⚠️  No URLs found for $group_name"
        return
    fi
    
    # Remove duplicates and empty lines
    unique_urls=$(echo "$all_urls" | sort -u | sed '/^$/d')
    
    # Count total unique URLs
    total=$(echo "$unique_urls" | wc -l)
    
    # Create JSON in the format required by /api/analysis/by-url
    echo "$unique_urls" | jq -R -s '{
        urls: split("\n") | map(select(length > 0)),
        model: "'"$MODEL"'"
    }' > "$OUTPUT_DIR/${group_name}.json"
    
    echo ""
    echo "  ✓ Total unique URLs: $total"
    echo "  ✓ Saved to: $OUTPUT_DIR/${group_name}.json"
}

# ==========================================
# 1. TRADIZIONE E SICUREZZA (Conservation)
# ==========================================
create_feed "conservation" \
    "difesa+confini" \
    "famiglia+tradizionale" \
    "identita+nazionale" \
    "sicurezza+pubblica" \
    "legge+ordine" \
    "sovranita+nazionale"

# ==========================================
# 2. POTERE E SUCCESSO (Self-Enhancement)
# ==========================================
create_feed "self_enhancement" \
    "mentalita+vincente" \
    "meritocrazia" \
    "leader+forte" \
    "successo+personale" \
    "ambizione" \
    "potere+autorita"

# ==========================================
# 3. UNIVERSALISMO E GIUSTIZIA (Self-Transcendence)
# ==========================================
create_feed "self_transcendence" \
    "emergenza+climatica" \
    "giustizia+sociale" \
    "diritti+umani" \
    "uguaglianza" \
    "nessuno+e+illegale" \
    "solidarieta"

# ==========================================
# 4. APERTURA E AUTONOMIA (Openness to Change)
# ==========================================
create_feed "openness_to_change" \
    "liberta+espressione" \
    "rompere+gli+schemi" \
    "pensiero+critico" \
    "nuove+esperienze" \
    "creativita" \
    "autonomia"

# ==========================================
# 5. TOPICS SPECIFICI (Hot topics)
# ==========================================
create_feed "topic_immigration" \
    "immigrazione+ice+trump" \
    "migranti+terra" \
    "profughi" \
    "accoglienza"

create_feed "topic_politics" \
    "elezioni+politiche" \
    "governo" \
    "parlamento" \
    "democrazia"

create_feed "topic_economy" \
    "economia+inflazione" \
    "prezzi+crescita" \
    "lavoro" \
    "disoccupazione"

echo ""
echo "=========================================="
echo "Download complete!"
echo "=========================================="
echo ""
echo "Feeds saved in: $OUTPUT_DIR"
echo ""
ls -lh "$OUTPUT_DIR"/*.json
echo ""
echo "To analyze a feed:"
echo "  curl -X POST http://localhost:8080/api/analysis/by-url \\"
echo "    -H 'Content-Type: application/json' \\"
echo "    -d @feeds/conservation.json | jq '.'"