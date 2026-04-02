package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bluesky-social/indigo/api/bsky"
)

// AlgoHandler - Firma funzione per un algoritmo di feed
type AlgoHandler func(ctx context.Context, limit int, cursor string) (*bsky.FeedGetFeedSkeleton_Output, error)

// Algos - Registro degli algoritmi disponibili
// Aggiungi il tuo feed qui con un nome unico
// Esempio URI risultante: at://did:plc:xxx/app.bsky.feed.generator/mio-feed
var Algos = map[string]AlgoHandler{
	"schwartz": feedStatico,
}

// =============================================================================
// CONFIGURAZIONE FEED STATICO
// =============================================================================
// Inserisci qui gli URI dei post che vuoi mostrare nel tuo feed
// Formato: at://did:plc:XXX/app.bsky.feed.post/YYY
// Puoi trovare gli URI usando l'API di Bluesky o strumenti come https://bsky.app
var postStatici = []string{
	"at://did:plc:t6ubj2wlhc34awzcymh3qpur/app.bsky.feed.post/3mi363eeufs2g",
}

// =============================================================================
// FINE CONFIGURAZIONE - Non modificare sotto se non necessario
// =============================================================================

// feedStatico - Restituisce i post dalla lista hardcoded
// Supporta paginazione con cursor (indice del post)
func feedStatico(ctx context.Context, limit int, cursor string) (*bsky.FeedGetFeedSkeleton_Output, error) {
	// Calcola indice di partenza dal cursor
	startIdx := 0
	if cursor != "" {
		idx, err := strconv.Atoi(cursor)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor: %w", err)
		}
		startIdx = idx
	}

	// Calcola indice di fine rispettando il limit
	endIdx := startIdx + limit
	if endIdx > len(postStatici) {
		endIdx = len(postStatici)
	}

	// Estrai slice di post per questa pagina
	postSlice := postStatici[startIdx:endIdx]

	// Converte in formato skeleton (solo URI)
	feed := make([]*bsky.FeedDefs_SkeletonFeedPost, len(postSlice))
	for i, uri := range postSlice {
		feed[i] = &bsky.FeedDefs_SkeletonFeedPost{
			Post: uri,
		}
	}

	// Genera cursor per pagina successiva
	// nil se siamo arrivati alla fine della lista
	var nextCursor *string
	if endIdx < len(postStatici) {
		c := strconv.Itoa(endIdx)
		nextCursor = &c
	}

	return &bsky.FeedGetFeedSkeleton_Output{
		Cursor: nextCursor,
		Feed:   feed,
	}, nil
}

// Aggiungi qui altre funzioni per feed aggiuntivi se necessario
// Esempio:
//
// var altriPost = []string{...}
//
// func feedAlternativo(ctx context.Context, limit int, cursor string) (*bsky.FeedGetFeedSkeleton_Output, error) {
//     ...stessa logica di feedStatico ma con altriPost...
// }
//
// Poi registralo in Algos:
// var Algos = map[string]AlgoHandler{
//     "mio-feed": feedStatico,
//     "altro-feed": feedAlternativo,
// }
