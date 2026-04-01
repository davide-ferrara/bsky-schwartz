package models

type SwartzValue struct {
	ID          string
	Name        string
	Cluster     string
	Weight      float64
	Description string
}

var SwartzValues = []SwartzValue{
	{"tolerance", "Tolleranza", "Self-Transcendence", 0, "Accettazione e comprensione di chi è diverso da sé"},
	{"nature", "Natura", "Self-Transcendence", 0, "Preservazione dell'ambiente naturale"},
	{"universalism", "Eguaglianza", "Self-Transcendence", 0, "Impegno per l'eguaglianza, la giustizia e la protezione di tutti"},
	{"caring", "Altruismo", "Self-Transcendence", 0, "Devozione verso le persone a cui si tiene"},
	{"dependability", "Responsabilità", "Self-Transcendence", 0, "Essere responsabili e leali verso gli altri"},

	{"humility", "Umiltà", "Conservation", 0, "Essere umili"},
	{"inter_conf", "Rispetto interpersonale", "Conservation", 0, "Evitare di turbare o ferire le altre persone"},
	{"rule_conf", "Rispetto delle leggi", "Conservation", 0, "Rispetto di regole, leggi e obblighi formali"},
	{"tradition", "Tradizione", "Conservation", 0, "Mantenere e preservare le tradizioni culturali, familiari o religiose"},
	{"societal_sec", "Sicurezza sociale", "Conservation", 0, "Sicurezza e stabilità nella società più ampia"},
	{"personal_sec", "Sicurezza personale", "Conservation", 0, "Sicurezza nell'ambiente immediato"},

	{"face", "Reputazione", "Self-Enhancement", 0, "Sicurezza e potere attraverso il mantenimento dell'immagine pubblica e l'evitare l'umiliazione"},
	{"resources", "Ricchezza", "Self-Enhancement", 0, "Controllo di risorse materiali e sociali"},
	{"dominance", "Potere", "Self-Enhancement", 0, "Influenza e diritto di comandare"},
	{"achievement", "Successo", "Self-Enhancement", 0, "Successo secondo gli standard sociali"},

	{"hedonism", "Piacere", "Openness to Change", 0, "Edonismo e gratificazione sensoriale"},
	{"stimulation", "Novità", "Openness to Change", 0, "Eccitazione, stimolazione e cambiamento"},
	{"sd_action", "Libertà personale", "Openness to Change", 0, "La libertà di determinare le proprie azioni"},
	{"sd_thought", "Indipendenza di pensiero", "Openness to Change", 0, "La libertà di coltivare le proprie idee e abilità"},
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
