# Data Analysis Tool

Benchmark tool for comparing GPT vs Minimax Schwartz value analysis on Bluesky posts.

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
| `--output` | Output PNG file path (default: comparison_chart.png) |
| `--diff` | Show difference chart instead of side-by-side comparison |
| `--clusters` | Show cluster-level differences instead of individual values |

## Output

By default, the tool generates a **side-by-side comparison chart** showing absolute scores:

- **Blue bars**: GPT scores
- **Orange bars**: Minimax scores

Use `--diff` flag to show a **difference chart** instead:

- **Green bars**: GPT scored higher than Minimax
- **Red bars**: Minimax scored higher than GPT

Each value is scored 0-6 according to Schwartz Theory of Basic Human Values.
