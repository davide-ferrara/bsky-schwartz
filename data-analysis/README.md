# Data Analysis Tool

Benchmark tool for comparing GPT vs Qwen Schwartz value analysis on Bluesky posts.

## Setup

```bash
uv sync
```

## Usage

### Prerequisites

Start the bsky-schwartz server:

```bash
cd ../ && make run
```

### Fetch posts by query

```bash
uv run python -m data_analysis.main \
  --base-url http://localhost:8080 \
  --query "politics" \
  --limit 10 \
  --output diff_chart.png
```

### With specific URIs

```bash
uv run python -m data_analysis.main \
  --base-url http://localhost:8080 \
  --uris "at://did:plc:xxx/app.bsky.feed.post/yyy,at://did:plc:zzz/app.bsky.feed.post/www" \
  --output diff_chart.png
```

### Cluster-level differences

```bash
uv run python -m data_analysis.main \
  --base-url http://localhost:8080 \
  --query "italy" \
  --limit 20 \
  --clusters \
  --output cluster_chart.png
```

## Options

| Option | Description |
|--------|-------------|
| `--base-url` | Base URL of the API (default: http://localhost:8080) |
| `--query` | Search query to get post URIs |
| `--limit` | Number of posts to fetch (default: 10) |
| `--uris` | Comma-separated list of post URIs |
| `--output` | Output PNG file path (default: diff_chart.png) |
| `--clusters` | Show cluster-level differences instead of individual values |

## Output

The tool generates a bar chart showing the difference between GPT and Qwen scores:

- **Green bars**: GPT scored higher than Qwen
- **Red bars**: Qwen scored higher than GPT

Each value is scored 0-6 according to Schwartz Theory of Basic Human Values.
