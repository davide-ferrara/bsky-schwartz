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
pkg/
├── scorer/          # Core scoring logic and AI integration
│   ├── scorer.go    # Score calculation and AI prompting
│   ├── types.go     # Data structures (FeedItem, SchwartzValues, Config)
├── bluesky/         # Bluesky Social API client
│   ├── client.go    # Query posts, threads, and social features

cmd/
└── server/          # HTTP API server
    └── main.go
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
| **Language**    | Go 1.21+                                          |
| **HTTP**        | Gin Web Framework                                   |
| **Data Stream** | Bluesky API (via `github.com/bluesky-social/indigo`) |
| **AI/ML**      | OpenRouter API (supports OpenAI, Anthropic, etc.)  |
| **Protocol**    | AT Protocol (Authenticated Transfer)                 |

---

## 5. API Endpoints

| Method | Endpoint                 | Description                              |
| ------ | ------------------------ | ---------------------------------------- |
| `GET`  | `/health`                | Health check                             |
| `GET`  | `/api/analysis`          | Search posts and analyze with Schwartz    |
| `GET`  | `/api/search`           | Search and return only post URIs         |
| `GET`  | `/api/analysis/by-uri`  | Get a single post by URI and analyze     |

---

## 6. Scoring System

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

## 7. Configuration

Weights are configured in `config.json` under the `weights` section, organized by
political orientation (e.g., `left`, `right`). Each weight can be positive or
negative, allowing the score to reflect alignment or opposition to certain values.

---

## 8. Build & Run Commands

- **Build:** `make build`
- **Run:** `make run`
- **Clean:** `make clean`
- **Test:** `make test`

Always use `make build` to compile instead of `go build` directly.

---
