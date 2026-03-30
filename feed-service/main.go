// Package main - Entry point del feed generator per Bluesky
// Versione semplificata per feed statico (senza database)
package main

import (
	"context"   // Gestione contesto per cancellazione graceful
	"fmt"       // Formattazione stringhe
	"log/slog"  // Logging strutturato
	"net/http"  // Server HTTP
	"os"        // Sistema operativo (exit, args)
	"os/signal" // Gestione segnali sistema (SIGINT, SIGTERM)
	"syscall"   // Costanti segnali
	"time"      // Timeouts e durate
)

// main - Punto di ingresso dell'applicazione
// Flusso: config → server HTTP
func main() {
	// Crea logger strutturato con output su console
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// Carica configurazione da variabili d'ambiente o file .env
	cfg, err := LoadConfig()
	if err != nil {
		logger.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	// Crea contesto che si cancella automaticamente con SIGINT (Ctrl+C) o SIGTERM
	// Permette shutdown graceful del server
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Crea server HTTP con tutti gli endpoint Bluesky
	// Per feed statico non serve database
	handler := NewServer(cfg, logger)
	addr := fmt.Sprintf("%s:%d", cfg.ListenHost, cfg.Port)
	srv := &http.Server{Addr: addr, Handler: handler}

	// Goroutine per shutdown graceful
	// Quando ctx viene cancellato (SIGINT/SIGTERM), chiude il server orderly
	go func() {
		<-ctx.Done()
		logger.Info("shutting down")

		// Timeout di 5 secondi per completare richieste in corso
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		srv.SetKeepAlivesEnabled(false) // Disabilita keep-alive
		srv.Shutdown(shutdownCtx)       // Chiude server gracefully
	}()

	// Avvia server HTTP (blocca fino a shutdown)
	logger.Info("running feed generator", "addr", fmt.Sprintf("http://%s", addr))
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("server error", "err", err)
		os.Exit(1)
	}
}
