# Bluesky Schwartz

AI-powered social media content analysis based on Schwartz Theory of Basic Human Values.

## Overview

Analizza post di Bluesky utilizzando i 19 valori fondamentali di Schwartz attraverso AI (OpenRouter).

## Flusso

```
Post URL → Fetch Bluesky → Build Prompt → AI Analysis → JSON Output
```

1. Fetch post da Bluesky (testo, link, metadata)
2. Costruisce prompt con dati post  
3. Chiama OpenRouter AI per analisi
4. Salva risultati in JSON

## Struttura

```
feed-generator/
├── main.go           # Entry point
├── client.go         # Bluesky client + estrazione dati
├── valueAnalysis.go  # AI analysis (OpenRouter)
├── prompts/
│   └── PROMPT_V3.md  # Prompt Schwartz values
└── feed/             # Feed URLs
```

## Setup

### Prerequisites

- Go 1.21+
- Bluesky account con app password
- OpenRouter API key

### Environment Variables

Creare `.env` in `feed-generator/`:

```bash
BSKY_HANDLE=tuo_handle.bsky.social
BSKY_APP_PASSWORD=tua_app_password
OPEN_ROUTER_KEY=tua_openrouter_key
```

## Usage

```bash
cd feed-generator
make build  # compila
make run    # esegue
```

## Output

`Posts_TIMESTAMP.json` con:

```json
{
  "AtURI": "at://did:plc:.../app.bsky.feed.post/...",
  "Text": "...",
  "AuthorName": "...",
  "Links": [...],
  "ValueAnalysis": {
    "Rating": {
      "Reputation": 0,
      "Power": 2,
      ...
    },
    "Reasoning": "...",
    "Stats": {
      "model": "openai/gpt-4o-mini",
      "response_time_ms": 2340,
      "cost_usd": 0.0012
    }
  }
}
```

## Schwartz Values (19)

| Cluster | Values |
|---------|--------|
| Openness to Change | Independent thoughts, Independent actions, Stimulation, Pleasure |
| Self-Enhancement | Achievement, Power, Wealth, Reputation |
| Conservation | Personal security, Societal security, Tradition, Lawfulness, Respect, Humility |
| Self-Transcendence | Caring, Responsibility, Equality, Nature, Tolerance |

## License

MIT