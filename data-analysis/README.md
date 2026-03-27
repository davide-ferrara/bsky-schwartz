# Data Analysis Tools

This directory contains tools for downloading and analyzing Bluesky feeds based on Schwartz values.

## Contents

- `download_feed.sh` - Download feeds from Bluesky API
- `analyze_feed.sh` - Analyze downloaded feeds with Schwartz values
- `feeds/` - Directory containing downloaded feed URLs
- `results/` - Directory containing analysis results

## Quick Start

### 1. Download Feeds

```bash
cd data-analysis
./download_feed.sh
```

This will create feed files organized by Schwartz values:
- `conservation.json` - Tradition & Security
- `self_enhancement.json` - Power & Success
- `self_transcendence.json` - Universalism & Justice
- `openness_to_change.json` - Freedom & Creativity
- `topic_*.json` - Current events feeds

### 2. Analyze Feeds

```bash
./analyze_feed.sh
```

This will analyze all downloaded feeds and save results in `results/`.

## Configuration

### Change Number of Results

Edit `download_feed.sh`:

```bash
LIMIT=50  # Change this to adjust results per query
```

### Change Model

Edit `download_feed.sh`:

```bash
MODEL="gpt4o"  # Available: gpt4o, qwen2, qwen3, gemini3, minimax2
```

## File Format

### Feed Format (for API)

```json
{
  "urls": [
    "https://bsky.app/profile/user.bsky.social/post/abc123",
    "https://bsky.app/profile/another.bsky.social/post/def456"
  ],
  "model": "gpt4o"
}
```

### Result Format

```json
[
  {
    "uri": "at://did:plc:xxx/app.bsky.feed.post/abc123",
    "text": "...",
    "values": {
      "achievement": 5,
      "caring": 3,
      ...
    },
    "score": -55.5,
    "model": "openai/gpt-4o-mini"
  }
]
```

## Use Cases

### Analyze Single Feed

```bash
curl -X POST http://localhost:8080/api/analysis/by-url \
  -H 'Content-Type: application/json' \
  -d @feeds/conservation.json | jq '.' > results/conservation_results.json
```

### Compare Scores

```bash
# Get all scores for conservation
jq '.[] | {uri: .uri, score: .score}' results/conservation_results.json

# Average score
jq '[.[] | .score] | add / length' results/conservation_results.json
```

### Filter by Score

```bash
# Get posts with score > 0
jq '.[] | select(.score > 0)' results/conservation_results.json
```

## Schwartz Values

The 19 Schwartz values mapped to feeds:

### Conservation (tradition & security)
- difesa confini, famiglia tradizionale, identità nazionale, sicurezza pubblica

### Self-Enhancement (power & success)
- mentalità vincente, meritocrazia, leader forte, successo personale

### Self-Transcendence (universalism & justice)
- emergenza climatica, giustizia sociale, diritti umani, uguaglianza

### Openness to Change (freedom & creativity)
- libertà espressione, rompere schemi, pensiero critico, nuove esperienze