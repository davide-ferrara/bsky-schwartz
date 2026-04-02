# Feed Generator

## Endpoint

### Feed Generator (http.ServeMux)

- `GET /.well-known/did.json` - Documento DID
- `GET /xrpc/_health` - Health check
- `GET /xrpc/app.bsky.feed.describeFeedGenerator` - Lista feed disponibili
- `GET /xrpc/app.bsky.feed.getFeedSkeleton` - Ottiene i post del feed

### Web App (Gin)

- `GET /` - Homepage con sliders
- `GET /login` - Pagina di login
- `POST /login` - Processa login
- `POST /logout` - Logout
- `GET /values` - Descrizione valori Schwartz
- `POST /preferences` - Salva preferenze
- `GET /static/*` - File statici
- `GET /lexicons/*` - Schemi AT Protocol

## Configurazione (.env)

```env
# Feed Generator
FEEDGEN_HOSTNAME=davideferrara.xyz
FEEDGEN_PUBLISHER_DID=did:plc:YOUR_DID_HERE
FEEDGEN_PORT=8080
FEEDGEN_LISTENHOST=0.0.0.0

# Web App
SESSION_KEY=your-secret-key-here

# Bluesky (opzionale)
BSKY_HANDLE=your-handle.bsky.social
BSKY_APP_PASSWORD=your-app-password
```

## Build & Run

```bash
# Build
go build -o bin/unified

# Run
./bin/unified

# O con Air (hot reload)
air
```

## Flusso Dati

1. **Login**: Utente → Web App → Autenticazione Bluesky → Salva weights in
   memoria condivisa
2. **Visualizzazione**: Web App → Legge weights da memoria condivisa
   (istantaneo)
3. **Modifica**: Utente → Web App → Salva in PDS + Aggiorna memoria condivisa
4. **Feed Generation**: Feed Generator → Legge weights da memoria condivisa →
   Personalizza feed

## Tecnologie

- **Go 1.25+** - Runtime
- **Gin** - Router HTTP per la Web App
- **Templ** - Template HTML type-safe
- **AT Protocol (indigo)** - Integrazione Bluesky
- **Memoria condivisa (sync.Map)** - Storage temporaneo ad alte prestazioni

## Licenza

MIT

