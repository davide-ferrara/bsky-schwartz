// Package main - Configurazione del feed generator
// Carica le variabili d'ambiente per configurare il server
package main

import (
	"fmt"     // Formattazione errori
	"os"      // Accesso variabili ambiente
	"strconv" // Conversione stringhe/interi

	"github.com/joho/godotenv" // Carica file .env
)

// Config - Struttura di configurazione del feed generator
// Tutti i valori possono essere impostati via variabili d'ambiente
type Config struct {
	// Porta su cui ascolta il server HTTP (default: 3000)
	Port int

	// Host su cui bindare (default: "localhost")
	// Usa "0.0.0.0" per accettare connessioni esterne
	ListenHost string

	// Hostname pubblico del server (es. "feed.example.com")
	// Usato per generare URL esterne e verifiche DID
	Hostname string

	// DID del servizio (identificativo decentralizzato)
	// Formato: "did:web:tuo-dominio.com" per dominio verificato
	// Oppure: "did:plc:xxx" per account Bluesky
	ServiceDID string

	// DID del publisher (chi pubblica il feed)
	// Solitamente il tuo DID personale o dell'account Bluesky
	PublisherDID string
}

// LoadConfig - Carica configurazione da variabili d'ambiente
// Legge file .env se presente, poi controlla le variabili d'ambiente
// I valori mancanti usano i default
//
// Variabili d'ambiente MINIME richieste:
//
//	FEEDGEN_HOSTNAME      → Hostname (obbligatorio per produzione!)
//	FEEDGEN_SERVICE_DID   → Service DID (default: did:web:{hostname})
//	FEEDGEN_PUBLISHER_DID → Il TUO DID Bluesky (obbligatorio!)
//
// Variabili opzionali:
//
//	FEEDGEN_PORT      → Port (default: 3000)
//	FEEDGEN_LISTENHOST → ListenHost (default: "localhost")
func LoadConfig() (*Config, error) {
	// Carica file .env se presente; ignora errore se mancante
	godotenv.Load()

	// Legge hostname (IMPORTANTE: deve essere il tuo dominio pubblico)
	hostname := envOrDefault("FEEDGEN_HOSTNAME", "localhost")

	// Genera ServiceDID se non specificato
	// did:web:example.com significa "possiedo example.com"
	serviceDID := envOrDefault("FEEDGEN_SERVICE_DID", fmt.Sprintf("did:web:%s", hostname))

	// Legge porta
	port, err := envIntOrDefault("FEEDGEN_PORT", 3000)
	if err != nil {
		return nil, fmt.Errorf("invalid FEEDGEN_PORT: %w", err)
	}

	return &Config{
		Port:         port,
		ListenHost:   envOrDefault("FEEDGEN_LISTENHOST", "localhost"),
		Hostname:     hostname,
		ServiceDID:   serviceDID,
		PublisherDID: envOrDefault("FEEDGEN_PUBLISHER_DID", "did:example:alice"),
	}, nil
}

// envOrDefault - Legge variabile d'ambiente o restituisce default
func envOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// envIntOrDefault - Legge variabile d'ambiente come intero o restituisce default
func envIntOrDefault(key string, defaultVal int) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal, nil
	}
	return strconv.Atoi(val)
}
