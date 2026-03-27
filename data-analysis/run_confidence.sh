#!/bin/bash

# Run confidence analysis on all Schwartz value feeds
# Uses multi-model comparison (gpt4o vs gemini3) by default

set -e

FEEDS=("conservation" "self_enhancement" "self_transcendence" "openness_to_change")
MODELS="${MODELS:-gpt4o,gemini3}"
RUNS="${RUNS:-5}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

if [ ! -d "results" ]; then
    mkdir -p results
fi

echo "========================================" >&2
echo "Confidence Analysis - All Feeds" >&2
echo "Models: $MODELS | Runs: $RUNS" >&2
echo "========================================" >&2

for feed in "${FEEDS[@]}"; do
    if [ -f "feeds/${feed}.json" ]; then
        echo "" >&2
        echo "Processing: ${feed}" >&2
        echo "----------------------------------------" >&2
        
        uv run python confidence_analysis.py \
            --feed "feeds/${feed}.json" \
            --models "$MODELS" \
            --runs "$RUNS" \
            --output "results/${feed}_confidence.json" \
            --chart "results/${feed}_chart.png"
    else
        echo "WARNING: Feed file not found: feeds/${feed}.json" >&2
    fi
done

echo "" >&2
echo "========================================" >&2
echo "All feeds processed successfully!" >&2
echo "Results saved in: results/*_confidence.json" >&2
echo "Charts saved in: results/*_chart.png" >&2
echo "========================================" >&2