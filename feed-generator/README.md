# Schwartz Values Feed Generator

Analizza post Bluesky con i valori Schwartz tramite AI.

## Flusso

```
Post URL → Fetch dati Bluesky → Build prompt → OpenRouter AI → JSON output
```

1. Fetch post da Bluesky (testo, link, immagini, metadata)
2. Costruisce prompt con dati post
3. Chiama OpenRouter AI per analisi valori Schwartz
4. Salva risultati in JSON

## Uso

```bash
make build# compila
make run   # esegue
```

## Configurazione

Creare `.env`:

```
BSKY_HANDLE=tuo_handle.bsky.social
BSKY_APP_PASSWORD=tua_app_password
OPEN_ROUTER_KEY=tua_openrouter_key
```

## Output

`Posts_TIMESTAMP.json` con analisi e statistiche AI.