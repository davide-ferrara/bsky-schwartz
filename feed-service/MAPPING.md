# Schwartz Values Mapping

## Mappatura Completa

I nomi delle chiavi JSON sono stati aggiornati per corrispondere ai nomi semplificati.

| JSON Key | Codice ID | Nome Italiano | Cluster |
|----------|-----------|---------------|---------|
| tolerance | tolerance | Tolleranza | Self-Transcendence |
| nature | nature | Natura | Self-Transcendence |
| equality | equality | Eguaglianza | Self-Transcendence |
| caring | caring | Altruismo | Self-Transcendence |
| responsibility | responsibility | Responsabilità | Self-Transcendence |
| humility | humility | Umiltà | Conservation |
| respect | respect | Rispetto interpersonale | Conservation |
| lawfulness | lawfulness | Rispetto delle leggi | Conservation |
| tradition | tradition | Tradizione | Conservation |
| societal_security | societal_security | Sicurezza sociale | Conservation |
| personal_security | personal_security | Sicurezza personale | Conservation |
| reputation | reputation | Reputazione | Self-Enhancement |
| wealth | wealth | Ricchezza | Self-Enhancement |
| power | power | Potere | Self-Enhancement |
| achievement | achievement | Successo | Self-Enhancement |
| pleasure | pleasure | Piacere | Openness to Change |
| stimulation | stimulation | Novità | Openness to Change |
| independent_actions | independent_actions | Libertà personale | Openness to Change |
| independent_thoughts | independent_thoughts | Indipendenza di pensiero | Openness to Change |

## Struttura JSON nel PDS

```json
{
  "$type": "com.schwartz.values",
  "weights": {
    "tolerance": 0.5,
    "nature": 0.0,
    "equality": -0.25,
    ...
  },
  "updatedAt": "2026-04-01T12:00:00Z"
}
```

## File Modificati

1. `lexicons/com/schwartz/values.json` - Schema AT Protocol
2. `internal/models/schwartz.go` - Definizione valori in italiano
3. `internal/models/weights.go` - Tipo mappa weights (invariato)

## Note

- I nomi italiani sono mantenuti per la UI
- I nomi inglesi semplificati sono usati come ID e nel JSON
- I cluster Schwartz sono preserved (Self-Transcendence, Conservation, Self-Enhancement, Openness to Change)