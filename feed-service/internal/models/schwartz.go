package models

type SwartzValue struct {
	ID          string
	Name        string
	Cluster     string
	Weight      float64
	Description string
}

var SwartzValues = []SwartzValue{
	// Self-Transcendence
	{"tolerance", "Tolleranza", "Self-Transcendence", 0, "Accettazione e comprensione di chi è diverso da sé"},
	{"nature", "Natura", "Self-Transcendence", 0, "Preservazione dell'ambiente naturale"},
	{"equality", "Eguaglianza", "Self-Transcendence", 0, "Impegno per l'eguaglianza, la giustizia e la protezione di tutti"},
	{"caring", "Altruismo", "Self-Transcendence", 0, "Devozione verso le persone a cui si tiene"},
	{"responsibility", "Responsabilità", "Self-Transcendence", 0, "Essere responsabili e leali verso gli altri"},

	// Conservation
	{"humility", "Umiltà", "Conservation", 0, "Essere umili"},
	{"respect", "Rispetto interpersonale", "Conservation", 0, "Evitare di turbare o ferire le altre persone"},
	{"lawfulness", "Rispetto delle leggi", "Conservation", 0, "Rispetto di regole, leggi e obblighi formali"},
	{"tradition", "Tradizione", "Conservation", 0, "Mantenere e preservare le tradizioni culturali, familiari o religiose"},
	{"societal_security", "Sicurezza sociale", "Conservation", 0, "Sicurezza e stabilità nella società più ampia"},
	{"personal_security", "Sicurezza personale", "Conservation", 0, "Sicurezza nell'ambiente immediato"},

	// Self-Enhancement
	{"reputation", "Reputazione", "Self-Enhancement", 0, "Sicurezza e potere attraverso il mantenimento dell'immagine pubblica e l'evitare l'umiliazione"},
	{"wealth", "Ricchezza", "Self-Enhancement", 0, "Controllo di risorse materiali e sociali"},
	{"power", "Potere", "Self-Enhancement", 0, "Influenza e diritto di comandare"},
	{"achievement", "Successo", "Self-Enhancement", 0, "Successo secondo gli standard sociali"},

	// Openness to Change
	{"pleasure", "Piacere", "Openness to Change", 0, "Edonismo e gratificazione sensoriale"},
	{"stimulation", "Novità", "Openness to Change", 0, "Eccitazione, stimolazione e cambiamento"},
	{"independent_actions", "Libertà personale", "Openness to Change", 0, "La libertà di determinare le proprie azioni"},
	{"independent_thoughts", "Indipendenza di pensiero", "Openness to Change", 0, "La libertà di coltivare le proprie idee e abilità"},
}

func GetValuesByCluster(cluster string) []SwartzValue {
	var result []SwartzValue
	for _, v := range SwartzValues {
		if v.Cluster == cluster {
			result = append(result, v)
		}
	}
	return result
}
