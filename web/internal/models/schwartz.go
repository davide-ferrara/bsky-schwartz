package models

type SwartzValue struct {
	ID          string
	Name        string
	Cluster     string
	ClusterName string
	Weight      float64
	Description string
}

var SwartzValues = []SwartzValue{
	{"tolerance", "Tolleranza", "Self-Transcendence", "Autotrascendenza", 0, "Accettazione e comprensione di chi è diverso da sé"},
	{"nature", "Natura", "Self-Transcendence", "Autotrascendenza", 0, "Preservazione dell'ambiente naturale"},
	{"universalism", "Eguaglianza", "Self-Transcendence", "Autotrascendenza", 0, "Impegno per l'eguaglianza, la giustizia e la protezione di tutti"},
	{"caring", "Altruismo", "Self-Transcendence", "Autotrascendenza", 0, "Devozione verso le persone a cui si tiene"},
	{"dependability", "Responsabilità", "Self-Transcendence", "Autotrascendenza", 0, "Essere responsabili e leali verso gli altri"},

	{"humility", "Umiltà", "Conservation", "Conservazione", 0, "Essere umili"},
	{"inter_conf", "Rispetto interpersonale", "Conservation", "Conservazione", 0, "Evitare di turbare o ferire le altre persone"},
	{"rule_conf", "Rispetto delle leggi", "Conservation", "Conservazione", 0, "Rispetto di regole, leggi e obblighi formali"},
	{"tradition", "Tradizione", "Conservation", "Conservazione", 0, "Mantenere e preservare le tradizioni culturali, familiari o religiose"},
	{"societal_sec", "Sicurezza sociale", "Conservation", "Conservazione", 0, "Sicurezza e stabilità nella società più ampia"},
	{"personal_sec", "Sicurezza personale", "Conservation", "Conservazione", 0, "Sicurezza nell'ambiente immediato"},

	{"face", "Reputazione", "Self-Enhancement", "Autoaffermazione", 0, "Sicurezza e potere attraverso il mantenimento dell'immagine pubblica e l'evitare l'umiliazione"},
	{"resources", "Ricchezza", "Self-Enhancement", "Autoaffermazione", 0, "Controllo di risorse materiali e sociali"},
	{"dominance", "Potere", "Self-Enhancement", "Autoaffermazione", 0, "Influenza e diritto di comandare"},
	{"achievement", "Successo", "Self-Enhancement", "Autoaffermazione", 0, "Successo secondo gli standard sociali"},

	{"hedonism", "Piacere", "Openness to Change", "Apertura al cambiamento", 0, "Edonismo e gratificazione sensoriale"},
	{"stimulation", "Novità", "Openness to Change", "Apertura al cambiamento", 0, "Eccitazione, stimolazione e cambiamento"},
	{"sd_action", "Libertà personale", "Openness to Change", "Apertura al cambiamento", 0, "La libertà di determinare le proprie azioni"},
	{"sd_thought", "Indipendenza di pensiero", "Openness to Change", "Apertura al cambiamento", 0, "La libertà di coltivare le proprie idee e abilità"},
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
