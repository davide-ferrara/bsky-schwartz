#!/bin/bash

# ==========================================
# Analyze downloaded feeds using /api/analysis/by-url
# ==========================================

set -e

API_URL="http://localhost:8080/api/analysis/by-url"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FEEDS_DIR="$SCRIPT_DIR/feeds"
RESULTS_DIR="$SCRIPT_DIR/results"

mkdir -p "$RESULTS_DIR"

analyze_feed() {
    local feed_file="$1"
    local output_name=$(basename "$feed_file" .json)
    
    echo "Analyzing: $output_name"
    
    if [ ! -f "$feed_file" ]; then
        echo "  ⚠️  File not found: $feed_file"
        return
    fi
    
    # Count URLs in feed
    url_count=$(jq '.urls | length' "$feed_file")
    echo "  Processing $url_count URLs..."
    
    # Send POST request to API
    response=$(curl -s -X POST "$API_URL" \
        -H "Content-Type: application/json" \
        -d @"$feed_file")
    
    # Check for errors
    if echo "$response" | jq -e '.error' > /dev/null 2>&1; then
        echo "  ❌ Error: $(echo "$response" | jq -r '.error')"
        return
    fi
    
    # Save results
    output_file="$RESULTS_DIR/${output_name}_results.json"
    echo "$response" | jq '.' > "$output_file"
    
    # Count results
    result_count=$(echo "$response" | jq 'length')
    echo "  ✓ Analyzed $result_count posts"
    echo "  ✓ Saved to: $output_file"
    echo ""
}

# ==========================================
# Analyze Schwartz value groups
# ==========================================
echo "Starting analysis of Schwartz value groups..."
echo ""

analyze_feed "$FEEDS_DIR/conservation.json"
analyze_feed "$FEEDS_DIR/self_enhancement.json"
analyze_feed "$FEEDS_DIR/self_transcendence.json"
analyze_feed "$FEEDS_DIR/openness_to_change.json"

# ==========================================
# Analyze Topical feeds
# ==========================================
echo "Analyzing topical feeds..."
echo ""

analyze_feed "$FEEDS_DIR/topic_immigration.json"
analyze_feed "$FEEDS_DIR/topic_politics.json"
analyze_feed "$FEEDS_DIR/topic_economy.json"

echo "=========================================="
echo "Analysis complete!"
echo "=========================================="
echo ""
echo "Results saved in: $RESULTS_DIR"
echo ""
echo "To view results:"
echo "  jq '.' $RESULTS_DIR/conservation_results.json"
echo ""
echo "To analyze scores:"
echo "  jq '.[] | {uri: .uri, score: .score}' $RESULTS_DIR/conservation_results.json"