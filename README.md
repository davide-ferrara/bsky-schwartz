# Bluesky Schwartz

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
[![Paper](https://img.shields.io/badge/Paper-arXiv:2509.14434-B31B1B?style=flat-square&logo=arxiv&logoColor=white)](https://arxiv.org/pdf/2509.14434)

**AI-Powered Social Media Content Analysis based on Schwartz Theory of Basic
Human Values**

bsky-schwartz is a Go application that analyzes Bluesky posts using the
[Schwartz Theory of Basic Human Values](https://en.wikipedia.org/wiki/Schwartz_theory_of_basic_human_values).
It leverages AI to evaluate content against 19 universal value dimensions and
calculates alignment scores.

## Background

This project is inspired by the paper
["Using Large Language Models to Assess Human Values in Social Media"](https://arxiv.org/abs/2509.14434)
which explores how LLMs can be used to analyze and quantify human values
expressed in social media content.

## Overview

The project fetches posts from Bluesky Social and uses AI (via OpenRouter) to
analyze how each post reflects the 19 basic human values defined by Shalom
Schwartz. Each post receives a score that indicates its alignment with different
value orientations.

## Features

- **Schwartz Value Analysis**: Evaluates posts against 19 universal human values
- **AI-Powered Scoring**: Uses OpenRouter-compatible AI models for analysis
- **Configurable Weights**: Supports different political/value orientations
  (left, right, etc.)
- **Bluesky Integration**: Direct integration with Bluesky Social API
- **REST API Ready**: Structured for easy HTTP API deployment
- **Thread-Safe Configuration**: Singleton config with sync.Once for safe
  concurrent access

## Architecture

```
pkg/
├── scorer/          # Core scoring logic and AI integration
│   ├── scorer.go    # Score calculation and AI prompting
│   ├── types.go     # Data structures (FeedItem, SchwartzValues, Config)
│   └── handler.go   # HTTP handlers (if applicable)
└── bluesky/         # Bluesky Social API client
    └── client.go    # Query posts, threads, and social features

cmd/
└── server/          # Main application entrypoint
    └── main.go

data/
├── SCHWARTZ.md       # Schwartz value definitions and scoring instructions
└── PROMPT.md        # AI prompt template
```

## The 19 Schwartz Values

| Value ID        | Cluster            | Description                               |
| --------------- | ------------------ | ----------------------------------------- |
| `sd_thought`    | Openness to Change | Freedom to cultivate one's own ideas      |
| `sd_action`     | Openness to Change | Freedom to determine one's own actions    |
| `stimulation`   | Openness to Change | Excitement and stimulation                |
| `hedonism`      | Openness to Change | Pleasure and sensuous gratification       |
| `achievement`   | Self-Enhancement   | Success according to social standards     |
| `dominance`     | Self-Enhancement   | Influence and power over others           |
| `resources`     | Self-Enhancement   | Control of material and social resources  |
| `face`          | Self-Enhancement   | Maintaining public image                  |
| `personal_sec`  | Conservation       | Personal safety and security              |
| `societal_sec`  | Conservation       | Safety and stability in society           |
| `tradition`     | Conservation       | Preserving cultural and religious customs |
| `rule_conf`     | Conservation       | Compliance with rules and laws            |
| `inter_conf`    | Conservation       | Respect for others and social norms       |
| `humility`      | Conservation       | Recognizing one's insignificance          |
| `caring`        | Self-Transcendence | Devotion to welfare of others             |
| `dependability` | Self-Transcendence | Reliability and loyalty                   |
| `universalism`  | Self-Transcendence | Justice and equality for all              |
| `nature`        | Self-Transcendence | Preservation of nature                    |
| `tolerance`     | Self-Transcendence | Acceptance of those different             |

## Setup

### Prerequisites

- Go 1.21+
- Bluesky account with app password
- OpenRouter API key

### Environment Variables

Create a `.env` file in the project root:

```bash
# Bluesky credentials
BSKY_HANDLE=your-handle.bsky.social
BSKY_APP_PASSWORD=your-app-password

# OpenRouter API key
OPEN_ROUTER_KEY=sk-or-v1-...

# Database (optional)
DB_HOST=localhost
DB_PORT=5433
DB_USER=admin
DB_PASSWORD=admin
DB_NAME=bsky

# Application
PORT=8080
CONFIG_PATH=config.json
```

### Configuration

Edit `config.json` to customize:

```json
{
  "models": {
    "gpt": "openai/gpt-4o-mini",
    "qwen": "qwen/qwen-2.5-72b-instruct"
  },
  "weights": {
    "left": {
      "Stimulation": -1,
      "Hedonism": -5,
      ...
    }
  },
  "ai": {
    "prompt": "data/PROMPT.md",
    "schwartz": "data/SCHWARTZ.md"
  }
}
```

## Usage

### Run the Server

```bash
make run
```

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

## API Reference

### Examples

```bash
# Health check
curl http://localhost:8080/health

# Search posts and analyze with Schwartz values
curl "http://localhost:8080/api/analysis?query=cats&limit=2&model=qwen"

# Search and return only post URIs
curl "http://localhost:8080/api/search?query=cats&limit=10"

# Get a single post by URI and analyze
curl "http://localhost:8080/api/analysis/by-uri?uri=at://did:plc:d5v2lwniz6g57usqe2fzxzgt/app.bsky.feed.post/3mhvowobwbc2g&model=qwen"
```

### Endpoints

| Method | Endpoint               | Description                            |
| ------ | ---------------------- | -------------------------------------- |
| `GET`  | `/health`              | Health check                           |
| `GET`  | `/api/analysis`        | Search posts and analyze with Schwartz |
| `GET`  | `/api/search`          | Search and return only post URIs       |
| `GET`  | `/api/analysis/by-uri` | Get a single post by URI and analyze   |

### Query Parameters

| Parameter | Type   | Default | Description                                       |
| --------- | ------ | ------- | ------------------------------------------------- |
| `query`   | string | -       | Search query                                      |
| `limit`   | int    | `10`    | Number of posts                                   |
| `model`   | string | `gpt`   | AI model to use                                   |
| `uri`     | string | -       | Post AT Protocol URI (for `/api/analysis/by-uri`) |

## Log Analysis with `jq`

The server writes structured JSON logs to `logs/server-{YYYY-MM-DD}.log`.  
You can use `jq` to query and analyze these logs:

### Read All Logs (Pretty Print)

```bash
jq . logs/server-2026-03-27.log
```

### Filter by Log Level

```bash
# Only errors
jq 'select(.level == "ERROR")' logs/server-2026-03-27.log

# Only warnings and errors
jq 'select(.level == "WARN" or .level == "ERROR")' logs/server-2026-03-27.log
```

### Filter by Request ID

Trace all operations for a specific request:

```bash
jq 'select(.request_id == "a1b2c3d4")' logs/server-2026-03-27.log
```

### View AI Analysis Metrics

```bash
# Show all AI analysis completions with costs and tokens
jq 'select(.msg == "ai analysis completed")' logs/server-2026-03-27.log

# Extract only cost and tokens
jq 'select(.msg == "ai analysis completed") | {time: .time, model: .model, tokens: .tokens_used, cost: .cost_usd}' logs/server-2026-03-27.log
```

### Calculate Total Costs

```bash
# Sum all AI costs for the day
jq 'select(.msg == "ai analysis completed") | .cost_usd' logs/server-2026-03-27.log | awk '{sum+=$1} END {print "Total: $"sum}'
```

### View Request Durations

```bash
# Show all completed requests with duration
jq 'select(.msg == "request completed") | {time: .time, path: .path, status: .status, duration_ms: .duration_ms}' logs/server-2026-03-27.log

# Find slowest requests (> 5 seconds)
jq 'select(.msg == "request completed" and .duration_ms > 5000)' logs/server-2026-03-27.log
```

### Track Bluesky API Performance

```bash
# Show all Bluesky search operations
jq 'select(.msg | contains("bluesky"))' logs/server-2026-03-27.log

# Average response time for Bluesky searches
jq 'select(.msg == "bluesky search completed") | .duration_ms' logs/server-2026-03-27.log | awk '{sum+=$1; count++} END {print "Avg: "sum/count" ms"}'
```

### View Errors with Context

```bash
# Show error messages with context
jq -C 'select(.level == "ERROR")' logs/server-2026-03-27.log | less -R
```

### Live Tail (Real-time Monitoring)

```bash
# Monitor logs in real-time
tail -f logs/server-$(date +%Y-%m-%d).log | jq .
```

### Aggregate Statistics

```bash
# Count requests by status code
jq -r 'select(.msg == "request completed") | .status' logs/server-2026-03-27.log | sort | uniq -c

# Count AI analyses per model
jq -r 'select(.msg == "ai analysis completed") | .model' logs/server-2026-03-27.log | sort | uniq -c
```

## Scoring System

Each value is scored 0-6:

- **0**: Value not present or contradicted
- **1-2**: Value slightly reflected
- **3-4**: Value moderately reflected
- **5-6**: Value strongly reflected

The final score is calculated as a weighted sum:

```
score = Σ(value_i × weight_i)
```

Negative weights apply a penalty multiplier (2.0x) to allow for differentiation
between opposing values.

## License

MIT
