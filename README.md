# Bluesky Schwartz

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
[![Paper](https://img.shields.io/badge/Paper-arXiv:2509.14434-B31B1B?style=flat-square&logo=arxiv&logoColor=white)](https://arxiv.org/pdf/2509.14434)

**AI-Powered Social Media Content Analysis based on Schwartz Theory of Basic Human Values**

bsky-schwartz is a Go application that analyzes Bluesky posts using the
[Schwartz Theory of Basic Human Values](https://en.wikipedia.org/wiki/Schwartz_theory_of_basic_human_values).
It leverages AI to evaluate content against 19 universal value dimensions and
calculates alignment scores.

## Background

This project is inspired by the paper ["Using Large Language Models to Assess
Human Values in Social Media"](https://arxiv.org/abs/2509.14434) which explores
how LLMs can be used to analyze and quantify human values expressed in social
media content.

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

### Endpoints

| Method | Endpoint  | Description  |
| ------ | --------- | ------------ |
| `GET`  | `/health` | Health check |

### Query Parameters

| Parameter | Type   | Default | Description     |
| --------- | ------ | ------- | --------------- |
| `query`   | string | `test`  | Search query    |
| `limit`   | int    | `10`    | Number of posts |
| `model`   | string | `gpt`   | AI model to use |

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
