// Package main - Tipid e memoria condivisa
package main

import "sync"

// User - Dati utente dalla sessione
type User struct {
	Handle      string
	AppPassword string
}

// UserWeights - Memoria condivisa per i weights degli utenti
// Chiave: DID (espresso come stringa)
// Valore: map[string]float64 (weights per ogni valore Schwartz)
type UserWeights struct {
	mu      sync.RWMutex
	weights map[string]map[string]float64 // map[DID]map[valueID]weight
}

// GlobalUserWeights - Istanza globale dei weights utente
var GlobalUserWeights = &UserWeights{
	weights: make(map[string]map[string]float64),
}

// Set - Imposta i weights per un utente (DID)
func (uw *UserWeights) Set(did string, weights map[string]float64) {
	uw.mu.Lock()
	defer uw.mu.Unlock()
	uw.weights[did] = weights
}

// Get - Ottiene i weights per un utente (DID)
// Restituisce nil se l'utente non ha weights
func (uw *UserWeights) Get(did string) map[string]float64 {
	uw.mu.RLock()
	defer uw.mu.RUnlock()
	if weights, ok := uw.weights[did]; ok {
		// Copia i weights per evitare race condition
		result := make(map[string]float64)
		for k, v := range weights {
			result[k] = v
		}
		return result
	}
	return nil
}
