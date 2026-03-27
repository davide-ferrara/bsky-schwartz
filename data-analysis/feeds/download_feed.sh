# ==========================================
# 1. TRADIZIONE E SICUREZZA (Conservation)
# ==========================================
# Fokus: Difesa dei confini e valori classici
curl "http://localhost:8080/api/search?query=difesa+confini&limit=20" | jq >>conservation.json
curl "http://localhost:8080/api/search?query=famiglia+naturale&limit=20" | jq >>conservation.json
curl "http://localhost:8080/api/search?query=identita+nazionale&limit=20" | jq >>conservation.json
curl "http://localhost:8080/api/search?query=fermare+invasione&limit=20" | jq >>conservation.json

# ==========================================
# 2. POTERE E DOMINIO (Self-Enhancement)
# ==========================================
# Fokus: Ambizione, gerarchia e successo individuale
curl "http://localhost:8080/api/search?query=mentalita+vincente&limit=20" | jq >>self_enhancement.json
curl "http://localhost:8080/api/search?query=meritocrazia&limit=20" | jq >>self_enhancement.json
curl "http://localhost:8080/api/search?query=successo+personale&limit=20" | jq >>self_enhancement.json
curl "http://localhost:8080/api/search?query=leader+forte&limit=20" | jq >>self_enhancement.json

# ==========================================
# 3. UNIVERSALISMO (Self-Transcendence)
# ==========================================
# Fokus: Diritti, clima e giustizia sociale
curl "http://localhost:8080/api/search?query=emergenza+climatica&limit=20" | jq >>universalism.json
curl "http://localhost:8080/api/search?query=patriarcato&limit=20" | jq >>universalism.json
curl "http://localhost:8080/api/search?query=antifascismo&limit=20" | jq >>universalism.json
curl "http://localhost:8080/api/search?query=nessuno+e+illegale&limit=20" | jq >>universalism.json

# ==========================================
# 4. APERTURA AL CAMBIAMENTO (Openness to Change)
# ==========================================
# Fokus: Creatività, libertà di espressione e nuove esperienze
curl "http://localhost:8080/api/search?query=liberta+espressione&limit=20" | jq >>openness_to_change.json
curl "http://localhost:8080/api/search?query=rompere+gli+schemi&limit=20" | jq >>openness_to_change.json
curl "http://localhost:8080/api/search?query=pensiero+critico&limit=20" | jq >>openness_to_change.json
curl "http://localhost:8080/api/search?query=nuove+esperienze&limit=20" | jq >>openness_to_change.json
