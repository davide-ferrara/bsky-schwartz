// Package main - Server HTTP per Bluesky Feed Generator
// Implementa gli endpoint richiesti dal protocollo AT (Bluesky)
package main

import (
	"encoding/json" // Codifica/decodifica JSON
	"fmt"           // Formattazione stringhe
	"log/slog"      // Logging strutturato
	"net/http"      // Server HTTP
	"strconv"       // Conversione stringhe/numeri
	"strings"       // Manipolazione stringhe
)

// NewServer - Crea il server HTTP con tutti gli endpoint Bluesky
// Endpoint esposti:
//   - /.well-known/did.json          → Documento DID per verifica identità
//   - /xrpc/_health                   → Health check
//   - /xrpc/app.bsky.feed.describeFeedGenerator → Lista feed disponibili
//   - /xrpc/app.bsky.feed.getFeedSkeleton      → Ottiene i post del feed
func NewServer(cfg *Config, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	// Endpoint "well-known" per risolvere il DID del feed generator
	// Bluesky richiede questo per verificare l'identità del server
	mux.HandleFunc("GET /.well-known/did.json", handleWellKnown(cfg))

	// Health check semplice per monitoraggio
	mux.HandleFunc("GET /xrpc/_health", handleHealth)

	// Descrive quali feed sono disponibili su questo server
	// Bluesky chiama questo per scoprire i feed dell'utente
	mux.HandleFunc("GET /xrpc/app.bsky.feed.describeFeedGenerator", handleDescribeFeedGenerator(cfg))

	// Endpoint principale: restituisce la lista di post per un feed
	// Questo è chiamato quando un utente visualizza il feed
	mux.HandleFunc("GET /xrpc/app.bsky.feed.getFeedSkeleton", handleGetFeedSkeleton(cfg, logger))

	return mux
}

// handleHealth - Endpoint di health check
// Risponde "OK" se il server è vivo
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// handleWellKnown - Gestisce /.well-known/did.json
// Restituisce il documento DID (Decentralized Identifier)
// Permette a Bluesky di verificare che questo server possiede il dominio
//
// Formula: ServiceDID deve finire con Hostname
// Esempio: did:web:example.com richiede hostname=example.com
func handleWellKnown(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verifica che il ServiceDID corrisponda all'hostname
		// did:web:example.com → hostname deve essere example.com
		if !strings.HasSuffix(cfg.ServiceDID, cfg.Hostname) {
			http.NotFound(w, r)
			return
		}
		// Documento DID nel formato standard W3C
		writeJSON(w, map[string]any{
			"@context": []string{"https://www.w3.org/ns/did/v1"},
			"id":       cfg.ServiceDID,
			"service": []map[string]any{
				{
					"id":              "#bsky_fg",
					"type":            "BskyFeedGenerator",
					"serviceEndpoint": fmt.Sprintf("https://%s", cfg.Hostname),
				},
			},
		})
	}
}

// handleDescribeFeedGenerator - Restituisce la lista dei feed disponibili
// Bluesky chiama questo endpoint per mostrare quali feed questo server offre
// Risponde con array di URI del tipo: at://did:xxx/app.bsky.feed.generator/nome-feed
func handleDescribeFeedGenerator(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Costruisce la lista di feed a partire dagli algoritmi registrati
		feeds := make([]map[string]string, 0, len(Algos))
		for shortname := range Algos {
			feeds = append(feeds, map[string]string{
				// URI nel formato AT Protocol
				"uri": fmt.Sprintf("at://%s/app.bsky.feed.generator/%s", cfg.PublisherDID, shortname),
			})
		}
		writeJSON(w, map[string]any{
			"did":   cfg.ServiceDID,
			"feeds": feeds,
		})
	}
}

// handleGetFeedSkeleton - Endpoint principale per ottenere i post del feed
// Restituisce una "skeleton" (struttura base) con soli URI dei post
// Bluesky poi hydrata i post completo (contenuto testo, immagini, etc.)
//
// Parametri query:
//   - feed: URI del feed richiesto (at://did:xxx/app.bsky.feed.generator/nome)
//   - limit: numero massimo di post (default 50, max 200)
//   - cursor: punto di ripresa per paginazione (indice)
func handleGetFeedSkeleton(cfg *Config, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Estrae parametro "feed" dalla query
		feedParam := r.URL.Query().Get("feed")
		if feedParam == "" {
			writeError(w, http.StatusBadRequest, "missing feed parameter")
			return
		}

		// Parsa l'URI del feed: at://did:xxx/app.bsky.feed.generator/rkey
		// Rimuove prefisso "at://" e divide in 3 parti
		parts := strings.Split(strings.TrimPrefix(feedParam, "at://"), "/")
		if len(parts) != 3 {
			writeError(w, http.StatusBadRequest, "invalid feed URI")
			return
		}
		did, collection, rkey := parts[0], parts[1], parts[2]

		// Verifica che il feed sia gestito da questo server
		if did != cfg.PublisherDID || collection != "app.bsky.feed.generator" {
			writeError(w, http.StatusBadRequest, "UnsupportedAlgorithm")
			return
		}

		// Cerca l'algoritmo registrato con questo nome (rkey)
		algo, ok := Algos[rkey]
		if !ok {
			writeError(w, http.StatusBadRequest, "UnsupportedAlgorithm")
			return
		}

		// Parsa parametro "limit" (default 50, max 200)
		limit := 50
		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
				limit = min(parsed, 200)
			}
		}
		// Cursor per paginazione (indice dell'ultimo post inviato)
		cursor := r.URL.Query().Get("cursor")

		// Esegue l'algoritmo per ottenere i post
		result, err := algo(r.Context(), limit, cursor)
		if err != nil {
			logger.Error("algo error", "algo", rkey, "err", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		// Restituisce il risultato come JSON
		writeJSON(w, result)
	}
}

// writeJSON - Helper per scrivere risposta JSON
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

// writeError - Helper per scrivere errore JSON con status code
func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
