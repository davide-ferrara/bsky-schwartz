# Project: Schwartz-Values Algorithmic Agent for Bluesky

## 1. Vision

The goal of this project is to implement **Algorithmic Sovereignty** on the AT
Protocol (Bluesky). Instead of a black-box algorithm designed for engagement,
this system provides a transparent, value-driven feed based on **Schwartz's
Theory of Basic Human Values**.

The agent analyzes Bluesky posts and scores them against 19 universal human
values using AI, then calculates alignment scores based on configurable weights.

---

## 2. Core Architecture

```
internal/
├── config.go    # Configuration and feeds loading
├── client.go    # Bluesky API client
├── scorer.go    # Scoring logic and AI integration
├── logger.go    # Structured logging
├── prompts.go   # Prompt caching

cmd/
└── feedgen/
    └── main.go   # CLI entry point

prompts/
├── SCHWARTZ.md      # Full Schwartz values definition
├── SCHWARTZ_LITE.md # Simplified Schwartz values
├── PROMPT.md        # v1 AI prompt
└── PROMPT_V2.md     # v2 AI prompt
```

---

## 3. The 19 Schwartz Values

| Value ID        | Cluster                | Description                               |
| --------------- | --------------------- | ----------------------------------------- |
| `sd_thought`    | Openness to Change    | Freedom to cultivate one's own ideas       |
| `sd_action`     | Openness to Change    | Freedom to determine one's own actions     |
| `stimulation`   | Openness to Change    | Excitement and stimulation                |
| `hedonism`      | Openness to Change    | Pleasure and sensuous gratification        |
| `achievement`   | Self-Enhancement      | Success according to social standards      |
| `dominance`     | Self-Enhancement      | Influence and power over others           |
| `resources`     | Self-Enhancement      | Control of material and social resources   |
| `face`          | Self-Enhancement      | Maintaining public image                   |
| `personal_sec`  | Conservation          | Personal safety and security               |
| `societal_sec`  | Conservation          | Safety and stability in society           |
| `tradition`     | Conservation          | Preserving cultural and religious customs  |
| `rule_conf`     | Conservation          | Compliance with rules and laws             |
| `inter_conf`    | Conservation          | Respect for others and social norms       |
| `humility`      | Conservation          | Recognizing one's insignificance           |
| `caring`        | Self-Transcendence    | Devotion to welfare of others            |
| `dependability`  | Self-Transcendence    | Reliability and loyalty                   |
| `universalism`  | Self-Transcendence    | Justice and equality for all               |
| `nature`        | Self-Transcendence    | Preservation of nature                    |
| `tolerance`     | Self-Transcendence    | Acceptance of those different             |

---

## 4. Technical Stack

| Component       | Technology                                          |
| :-------------- | :-------------------------------------------------- |
| **Language**    | Go 1.26+                                          |
| **CLI**         | Native Go CLI                                       |
| **Bluesky API** | `github.com/bluesky-social/indigo`                  |
| **AI/ML**       | OpenRouter API (supports OpenAI, Anthropic, etc.)  |

---

## 5. Scoring System

Each value is scored 0-6:

- **0**: Value not present or contradicted
- **1-2**: Value slightly reflected
- **3-4**: Value moderately reflected
- **5-6**: Value strongly reflected

The final score is calculated as a weighted sum:

```
score = Σ(value_i × weight_i)
```

---

## 6. Configuration

- `config.json` - Models, weights, and AI prompt configuration
- `feeds.json` - List of feed URLs to process
- `.env` - Environment variables (BSKY_HANDLE, BSKY_APP_PASSWORD, OPEN_ROUTER_KEY)

---

## 7. Usage

```bash
# Run with defaults
make run

# Run with custom options
./bin/feedgen -config config.json -feeds feeds.json -model gpt -log info

# Build binary
make build
```

### CLI Flags

| Flag      | Default       | Description                    |
| --------- | ------------- | ------------------------------ |
| `-config` | `config.json` | Path to config file            |
| `-feeds`  | `feeds.json`  | Path to feeds file             |
| `-model`  | `gpt`         | Model key from config          |
| `-log`    | `info`        | Log level (debug/info/warn/error) |

### Environment Variables

Required in `.env`:
- `BSKY_HANDLE` - Bluesky handle
- `BSKY_APP_PASSWORD` - Bluesky app password
- `OPEN_ROUTER_KEY` - OpenRouter API key

---

## 8. Build & Run Commands

- **Build:** `make build`
- **Run:** `make run`
- **Clean:** `make clean`
- **Test:** `make test`

Always use `make build` to compile instead of `go build` directly.