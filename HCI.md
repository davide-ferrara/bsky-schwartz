# Disegno sperimentale

L'esperimento within-subject (a misure ripetute) è composto nel seguente modo.
Ogni partecipante è esposto a due condizioni:

- Feed A (baseline): ordinamento basato su logiche di engagement
- Feed B (value-based): ordinamento generato dall’algoritmo di value alignment

Entrambi i feed contengono lo stesso insieme di post e differiscono
esclusivamente per l’ordinamento, così da isolare l’effetto del ranking.

## Procedura

Lo studio si articola in tre fasi:

1. I partecipanti impostano i propri pesi sui valori tramite un’interfaccia web.
2. Vengono mostrati simultaneamente due feed (A e B), affiancati e con ordine
   randomizzato.
3. I partecipanti compilano un questionario di valutazione.

## Test 1: Allineamento percepito

Il test misura l’allineamento percepito tra i contenuti mostrati e i valori
dell’utente, seguendo un approccio user-centered ispirato alla letteratura sulla
valutazione dei recommender system (Pu et al., 2011; Knijnenburg et al., 2012).

Per ciascun feed, i partecipanti valutano le seguenti affermazioni su una scala
Likert da 1 (totalmente in disaccordo) a 7 (totalmente d’accordo):

- Questo feed riflette i miei valori personali.
- I contenuti che vedo sono coerenti con ciò che ritengo importante.
- Mi riconosco nei contenuti mostrati in questo feed.
- Questo feed rappresenta ciò a cui tengo.
- L’ordine dei contenuti rispecchia le mie priorità.

Per ogni partecipante viene calcolata la media delle risposte per Feed A e Feed
B. Le due condizioni vengono confrontate tramite test statistico per campioni
appaiati (t-test o Wilcoxon).

## Test 2: Preferenza

Il test di preferenza segue un approccio comparativo tipico degli studi HCI e
della valutazione user-centric dei recommender system (Pu et al., 2011;
Knijnenburg et al., 2012).

Ai partecipanti vengono mostrati estratti dei due feed (sottoinsiemi casuali
degli stessi post ordinati secondo le due condizioni) e viene chiesto:

- Quale dei due feed preferisci complessivamente? (Feed A / Feed B)
- Quale feed utilizzeresti in una piattaforma reale? (Feed A / Feed B)

In una variante più controllata, viene mostrato lo stesso insieme di post con
due ordinamenti diversi e viene chiesto:

- Quale ordinamento rappresenta meglio i tuoi interessi e valori?

Le preferenze vengono analizzate calcolando la distribuzione delle risposte e
applicando un test binomiale o chi-quadro.

## Test 3: Riconoscibilità

Il test misura la capacità degli utenti di distinguere tra diversi criteri di
ordinamento, in linea con studi sulla percezione degli algoritmi nei recommender
system (Ekstrand et al., 2014).

Vengono mostrati entrambi i feed ordinati secondo un valore specifico.

Domanda:

- Quale feed riflette meglio il valore X?

## Domanda aperta (opzionale)

- Perché hai preferito questo feed?

Le risposte vengono analizzate qualitativamente per supportare l’interpretazione
dei risultati quantitativi.
