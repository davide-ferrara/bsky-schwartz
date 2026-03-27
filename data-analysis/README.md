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

## Confidence Analysis

### Overview

Run statistical analysis with confidence intervals by analyzing each feed multiple times. This calculates mean, standard deviation, and 95% confidence intervals for all Schwartz values and scores.

**Two modes:**
- **Single-model**: Run multiple times with one model (baseline)
- **Multi-model comparison**: Compare multiple models side-by-side (default)

### Available Models

Add new models by editing `MODELS` dict in `confidence_analysis.py`:

```python
MODELS = {
    "gpt4o": "openai/gpt-4o-mini",
    "gemini3": "google/gemini-3.1-flash-lite-preview",
    # Add more models here
}
```

Default comparison: `gpt4o` vs `gemini3`

### Single-Model Analysis

Run analysis with a single model (backward compatible):

```bash
# Run 5 times with gpt4o
uv run python confidence_analysis.py --feed feeds/conservation.json --model gpt4o

# Custom number of runs
uv run python confidence_analysis.py --feed feeds/conservation.json --model gpt4o --runs 10

# Specify output paths
uv run python confidence_analysis.py --feed feeds/conservation.json \
    --model gpt4o \
    --output results/my_confidence.json \
    --chart results/my_chart.png
```

### Multi-Model Comparison

Compare multiple models on the same feed:

```bash
# Default: gpt4o vs gemini3
uv run python confidence_analysis.py --feed feeds/conservation.json

# Specify models to compare
uv run python confidence_analysis.py --feed feeds/conservation.json --models gpt4o,gemini3

# Future: compare 3+ models
uv run python confidence_analysis.py --feed feeds/conservation.json --models gpt4o,gemini3,qwen2

# Skip chart generation
uv run python confidence_analysis.py --feed feeds/conservation.json --no-chart
```

### Batch Analysis (All Feeds)

```bash
# Run confidence analysis on all feeds with default models
uv run ./run_confidence.sh
```

### Output Files

**JSON Output** (`results/{feed}_confidence.json`):

Single model:
```json
{
  "feed": "conservation",
  "runs": 5,
  "models": {
    "gpt4o": {
      "tradition": {"mean": 3.8, "std": 0.7, "ci_low": 3.1, "ci_high": 4.5},
      "caring": {"mean": 4.1, "std": 0.6, "ci_low": 3.4, "ci_high": 4.8},
      ...
      "score": {"mean": -12.3, "std": 2.1, "ci_low": -15.0, "ci_high": -9.6},
      "score_penalized": {"mean": -24.6, "std": 4.2, "ci_low": -30.0, "ci_high": -19.2}
    }
  },
  "timestamp": "2026-03-27T21:45:00"
}
```

Multi-model comparison:
```json
{
  "feed": "conservation",
  "runs": 5,
  "models": {
    "gpt4o": {
      "tradition": {"mean": 3.8, "std": 0.7, "ci_low": 3.1, "ci_high": 4.5},
      ...
    },
    "gemini3": {
      "tradition": {"mean": 3.5, "std": 0.9, "ci_low": 2.6, "ci_high": 4.4},
      ...
    }
  },
  "timestamp": "2026-03-27T21:45:00"
}
```

**Chart Output** (`results/{feed}_chart.png`):

Horizontal bar chart with two subplots:
- **Left**: Schwartz values (0-6 scale) with error bars showing 95% CI
- **Right**: Overall scores with error bars

**Colors**:
- gpt4o: Blue (#3B82F6)
- gemini3: Orange (#F97316)

**Console Output**:
```
============================================================
Confidence Analysis: conservation
Models: gpt4o, gemini3 | Runs: 5
============================================================

Processing model: gpt4o
  Run 1/5... ✓
  Run 2/5... ✓
  Run 3/5... ✓
  Run 4/5... ✓
  Run 5/5... ✓
  Calculating statistics...
  Done!

Processing model: gemini3
  Run 1/5... ✓
  Run 2/5... ✓
  Run 3/5... ✓
  Run 4/5... ✓
  Run 5/5... ✓
  Calculating statistics...
  Done!

Summary for conservation
============================================================

Openness to Change:
  sd_thought          :
    gpt4o           :  2.30 ± 0.50 (CI:  1.80 -  2.80)
    gemini3         :  2.40 ± 0.60 (CI:  1.80 -  3.00)
  ...

Self-Enhancement:
  ...

Conservation:
  tradition           :
    gpt4o           :  3.80 ± 0.70 (CI:  3.10 -  4.50)
    gemini3         :  3.50 ± 0.90 (CI:  2.60 -  4.40)
  ...

Self-Transcendence:
  ...

Overall Scores:
  score               :
    gpt4o           : -12.30 ± 2.10 (CI: -15.00 - -9.60)
    gemini3         : -10.50 ± 3.20 (CI: -14.30 - -6.70)
  score_penalized     :
    gpt4o           : -24.60 ± 4.20 (CI: -30.00 - -19.20)
    gemini3         : -21.00 ± 6.40 (CI: -28.60 - -13.40)
============================================================
```

### Statistical Methodology

- **Confidence Level**: 95%
- **Method**: Student's t-distribution (appropriate for small sample sizes)
- **Formula**: CI = mean ± t(n-1, 0.975) × (std / √n)
- **Data Points**: All values from all posts across all runs (e.g., 5 runs × 19 posts = 95 data points)
- **Error Bars**: Show confidence interval range around the mean

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